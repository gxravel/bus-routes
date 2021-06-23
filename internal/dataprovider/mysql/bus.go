package mysql

import (
	"context"

	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/logger"
	"github.com/gxravel/bus-routes/internal/model"
	"github.com/jmoiron/sqlx"
)

type BusStore struct {
	db        sqlx.ExtContext
	txer      dataprovider.Txer
	tableName string
}

func NewBusStore(db sqlx.ExtContext, txer dataprovider.Txer) *BusStore {
	return &BusStore{
		db:        db,
		txer:      txer,
		tableName: "bus",
	}
}

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
		eq["city"] = f.Cities
	}
	if len(f.Nums) > 0 {
		eq["num"] = f.Nums
	}

	return cond
}

func (s *BusStore) ByFilter(ctx context.Context, filter *dataprovider.BusFilter) (*model.Bus, error) {
	buses, err := s.ListByFilter(ctx, filter)

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

func (s *BusStore) ListByFilter(ctx context.Context, filter *dataprovider.BusFilter) ([]*model.Bus, error) {
	qb := sq.
		Select(
			"bus.id",
			"city.name as city",
			"num",
		).
		From(s.tableName).
		InnerJoin("city on bus.city_id=city.id").
		Where(busCond(filter))

	if filter.Paginator != nil {
		qb = withPaginator(qb, filter.Paginator)
	}

	result, err := selectContext(ctx, qb, s.tableName, s.db, TypeBus)
	if err != nil {
		return nil, err
	}
	return result.([]*model.Bus), nil
}

func (s *BusStore) Add(ctx context.Context, buses ...*model.Bus) error {
	var ids = make(map[string]int, len(buses))
	for _, bus := range buses {
		ids[bus.City] = -1
	}

	f := func(tx *dataprovider.Tx) error {
		err := CitiesIDs(ctx, ids, s.db, s.txer, tx)
		if err != nil {
			return err
		}
		qb := sq.Insert("bus").Columns("city_id", "num")
		for _, bus := range buses {
			id := ids[bus.City]
			if id < 0 {
				logger.FromContext(ctx).Debugf("bus [%s, %s] skipped", bus.City, bus.Num)
				continue
			}
			qb = qb.Values(id, bus.Num)
		}

		query, args, _, err := toSql(ctx, qb, s.tableName)
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return errors.Wrapf(err, "inserting buses with query %s", query)
		}
		return nil
	}

	return dataprovider.BeginAutoCommitedTx(ctx, s.txer, f)
}
