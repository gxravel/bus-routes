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

// BusStore is bus mysql store.
type BusStore struct {
	db        sqlx.ExtContext
	txer      dataprovider.Txer
	tableName string
}

// NewBusStore creates new instance of BusStore.
func NewBusStore(db sqlx.ExtContext, txer dataprovider.Txer) *BusStore {
	return &BusStore{
		db:        db,
		txer:      txer,
		tableName: "bus",
	}
}

// WithTx sets transaction as active connection.
func (s *BusStore) WithTx(tx *dataprovider.Tx) dataprovider.BusStore {
	return &BusStore{
		db:        tx,
		tableName: s.tableName,
	}
}

func busCond(f *dataprovider.BusFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	var cond sq.Sqlizer = eq

	if len(f.IDs) > 0 {
		eq["bus.id"] = f.IDs
	}
	if len(f.Cities) > 0 {
		eq["name"] = f.Cities
	}
	if len(f.Nums) > 0 {
		eq["num"] = f.Nums
	}

	return cond
}

func (s *BusStore) columns(filter *dataprovider.BusFilter) []string {
	var result = []string{
		"bus.id",
		"city.name as name",
		"num",
	}
	if filter.DoPreferIDs {
		result[1] = "bus.city_id"
	}
	return result
}

func (s *BusStore) joins(qb sq.SelectBuilder, filter *dataprovider.BusFilter) sq.SelectBuilder {
	if !filter.DoPreferIDs {
		qb = qb.Join("city ON bus.city_id = city.id")
	}
	return qb
}

// GetByFilter returns bus depend on received filters.
func (s *BusStore) GetByFilter(ctx context.Context, filter *dataprovider.BusFilter) (*model.Bus, error) {
	buses, err := s.GetListByFilter(ctx, filter)

	switch {
	case err != nil:
		return nil, err
	case len(buses) == 0:
		return nil, nil
	case len(buses) == 1:
		return buses[0], nil
	default:
		return nil, errors.New("fetched more than 1 bus")
	}
}

// GetListByFilter returns buses depend on received filters.
func (s *BusStore) GetListByFilter(ctx context.Context, filter *dataprovider.BusFilter) ([]*model.Bus, error) {
	qb := sq.
		Select(s.columns(filter)...).
		From(s.tableName).
		Where(busCond(filter))

	qb = s.joins(qb, filter)

	if filter.Paginator != nil {
		qb = withPaginator(qb, filter.Paginator)
	}

	result, err := selectContext(ctx, qb, s.tableName, s.db, TypeBus)
	if err != nil {
		return nil, err
	}
	return result.([]*model.Bus), nil
}

// Add creates new buses skipping those of with wrong city.
func (s *BusStore) Add(ctx context.Context, buses ...*model.Bus) error {
	var ids = make(map[string]int, len(buses))
	for _, bus := range buses {
		ids[bus.City] = 0
	}

	f := func(tx *dataprovider.Tx) error {
		err := CitiesIDs(ctx, ids, s.db, s.txer, tx)
		if err != nil {
			return err
		}
		qb := sq.Insert("bus").Columns("city_id", "num")
		for _, bus := range buses {
			id := ids[bus.City]
			if id == 0 {
				logger.FromContext(ctx).Debugf("bus [%s, %s] skipped", bus.City, bus.Number)
				continue
			}
			qb = qb.Values(id, bus.Number)
		}

		query, args, _, err := toSql(ctx, qb, s.tableName)
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return errors.Wrapf(err, "insert buses with query %s", query)
		}
		return nil
	}

	return dataprovider.BeginAutoCommitedTx(ctx, s.txer, f)
}
