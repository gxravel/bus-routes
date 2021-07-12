package mysql

import (
	"context"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	log "github.com/gxravel/bus-routes/internal/logger"
	"github.com/gxravel/bus-routes/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// StopStore is stop mysql store.
type StopStore struct {
	db        sqlx.ExtContext
	txer      dataprovider.Txer
	tableName string
}

// NewStopStore creates new instance of StopStore.
func NewStopStore(db sqlx.ExtContext, txer dataprovider.Txer) *StopStore {
	return &StopStore{
		db:        db,
		txer:      txer,
		tableName: "stop",
	}
}

// WithTx sets transaction as active connection.
func (s *StopStore) WithTx(tx *dataprovider.Tx) dataprovider.StopStore {
	return &StopStore{
		db:        tx,
		tableName: s.tableName,
	}
}

func stopCond(f *dataprovider.StopFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	var cond sq.Sqlizer = eq

	if len(f.IDs) > 0 {
		eq["stop.id"] = f.IDs
	}
	if len(f.Cities) > 0 {
		eq["city"] = f.Cities
	}
	if len(f.Addresses) > 0 {
		eq["address"] = f.Addresses
	}

	return cond
}

func (s *StopStore) columns(filter *dataprovider.StopFilter) []string {
	var result = []string{
		"stop.id",
		"city.name as city",
		"address",
	}

	if filter.DoPreferIDs {
		result[1] = "stop.city_id"
	}

	return result
}

func (s *StopStore) joins(qb sq.SelectBuilder, filter *dataprovider.StopFilter) sq.SelectBuilder {
	if !filter.DoPreferIDs {
		qb = qb.Join("city ON stop.city_id = city.id")
	}

	return qb
}

// GetByFilter returns stop depend on received filters.
func (s *StopStore) GetByFilter(ctx context.Context, filter *dataprovider.StopFilter) (*model.Stop, error) {
	stops, err := s.GetListByFilter(ctx, filter)

	switch {
	case err != nil:
		return nil, err
	case len(stops) == 0:
		return nil, nil
	case len(stops) == 1:
		return stops[0], nil
	default:
		return nil, errors.New("fetched more than 1 Stop")
	}
}

// GetListByFilter returns stops depend on received filters.
func (s *StopStore) GetListByFilter(ctx context.Context, filter *dataprovider.StopFilter) ([]*model.Stop, error) {
	qb := sq.
		Select(s.columns(filter)...).
		From(s.tableName).
		Where(stopCond(filter))

	qb = s.joins(qb, filter)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	message := "select " + s.tableName + " by filter with query " + query

	var result = make([]*model.Stop, 0)
	if err := sqlx.SelectContext(ctx, s.db, &result, query, args...); err != nil {
		return nil, errors.Wrapf(err, message)
	}

	return result, nil
}

// Add creates new stops skipping those of with wrong city.
func (s *StopStore) Add(ctx context.Context, stops ...*model.Stop) error {
	var ids = make(map[string]int, len(stops))
	for _, stop := range stops {
		ids[stop.City] = -1
	}

	f := func(tx *dataprovider.Tx) error {
		if err := getCitiesIDs(ctx, ids, s.db, s.txer, tx); err != nil {
			return err
		}

		qb := sq.Insert("stop").Columns("city_id", "address")
		for _, stop := range stops {
			id := ids[stop.City]
			if id < 0 {
				log.FromContext(ctx).
					Debugf("stop [%s, %s] skipped", stop.City, stop.Address)
				continue
			}
			qb = qb.Values(id, stop.Address)
		}

		if err := execContext(ctx, qb, s.tableName, tx); err != nil {
			return err
		}

		return nil
	}

	return dataprovider.BeginAutoCommitedTx(ctx, s.txer, f)
}

// Update updates stop's city_id and address.
func (s *StopStore) Update(ctx context.Context, stop *model.Stop) error {
	f := func(tx *dataprovider.Tx) error {
		id, err := getCityID(ctx, stop.City, s.db, s.txer, tx)
		if err != nil {
			return err
		}
		if id == 0 {
			err := errors.Errorf("did not found the city %s", stop.City)
			log.FromContext(ctx).Debug(err.Error())

			return err
		}

		qb := sq.Update(s.tableName).
			SetMap(map[string]interface{}{
				"city_id": id,
				"address": stop.Address,
			}).
			Where(sq.Eq{"id": stop.ID})

		if err := execContext(ctx, qb, s.tableName, s.db); err != nil {
			return err
		}

		return nil
	}

	return dataprovider.BeginAutoCommitedTx(ctx, s.txer, f)
}

// Delete deletes stop depend on received filter.
func (s *StopStore) Delete(ctx context.Context, filter *dataprovider.StopFilter) error {
	qb := sq.Delete(s.tableName).Where(stopCond(filter))
	return execContext(ctx, qb, s.tableName, s.db)
}
