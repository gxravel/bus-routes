package mysql

import (
	"context"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/logger"
	"github.com/gxravel/bus-routes/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type StopStore struct {
	db        sqlx.ExtContext
	txer      dataprovider.Txer
	tableName string
}

func NewStopStore(db sqlx.ExtContext, txer dataprovider.Txer) *StopStore {
	return &StopStore{
		db:        db,
		txer:      txer,
		tableName: "stop",
	}
}

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

func (s *StopStore) ByFilter(ctx context.Context, filter *dataprovider.StopFilter) (*model.Stop, error) {
	stops, err := s.ListByFilter(ctx, filter)

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

func (s *StopStore) ListByFilter(ctx context.Context, filter *dataprovider.StopFilter) ([]*model.Stop, error) {
	qb := sq.
		Select(s.columns(filter)...).
		From(s.tableName).
		Where(stopCond(filter))

	qb = s.joins(qb, filter)

	result, err := selectContext(ctx, qb, s.tableName, s.db, TypeStop)
	if err != nil {
		return nil, err
	}
	return result.([]*model.Stop), nil
}

func (s *StopStore) Add(ctx context.Context, stops ...*model.Stop) error {
	var ids = make(map[string]int, len(stops))
	for _, stop := range stops {
		ids[stop.City] = -1
	}

	f := func(tx *dataprovider.Tx) error {
		err := CitiesIDs(ctx, ids, s.db, s.txer, tx)
		if err != nil {
			return err
		}
		qb := sq.Insert("stop").Columns("city_id", "address")
		for _, stop := range stops {
			id := ids[stop.City]
			if id < 0 {
				logger.FromContext(ctx).Debugf("stop [%s, %s] skipped", stop.City, stop.Address)
				continue
			}
			qb = qb.Values(id, stop.Address)
		}

		query, args, _, err := toSql(ctx, qb, s.tableName)
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return errors.Wrapf(err, "inserting stops with query %s", query)
		}
		return nil
	}

	return dataprovider.BeginAutoCommitedTx(ctx, s.txer, f)
}

func (s *StopStore) Update(ctx context.Context, stop *model.Stop) error {
	f := func(tx *dataprovider.Tx) error {
		id, err := CityID(ctx, stop.City, s.db, s.txer, tx)
		if err != nil {
			return err
		}
		if id == 0 {
			err := errors.Errorf("did not found the city %s", stop.City)
			logger.FromContext(ctx).Debug(err.Error())
			return err
		}
		qb := sq.Update(s.tableName).Set("city_id", id).Set("Address", stop.Address).Where(sq.Eq{"id": stop.ID})
		query, args, _, err := toSql(ctx, qb, s.tableName)
		if err != nil {
			return err
		}

		result, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return errors.Wrapf(err, "updating stop with query %s", query)
		}
		num, err := result.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "failed to call RowsAffected")
		}
		if num == 0 {
			return errors.New("no rows affected: wrong id")
		}
		return nil
	}

	return dataprovider.BeginAutoCommitedTx(ctx, s.txer, f)
}

func (s *StopStore) Delete(ctx context.Context, filter *dataprovider.StopFilter) error {
	qb := sq.Delete(s.tableName).Where(stopCond(filter))
	return execContext(ctx, qb, s.tableName, s.txer)
}
