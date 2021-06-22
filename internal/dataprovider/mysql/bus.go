package mysql

import (
	"context"
	"database/sql"

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

func (s *BusStore) logger(ctx context.Context) logger.Logger {
	return logger.FromContext(ctx).WithStr("module", "mysql:bus")
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

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "creating sql query for getting buses by filter")
	}

	s.logger(ctx).
		WithFields(
			"query", query,
			"args", args).
		Debug("selecting buses by filter query SQL")

	buses := make([]*model.Bus, 0)
	if err = sqlx.SelectContext(ctx, s.db, &buses, query, args...); err != nil {
		return nil, errors.Wrapf(err, "selecting buses by filter from database with query %s", query)
	}

	return buses, nil
}

func (s *BusStore) New(ctx context.Context, buses ...*model.Bus) error {
	var ids = make(map[string]int, len(buses))
	for _, bus := range buses {
		ids[bus.City] = -1
	}

	var names = make([]string, 0, len(ids))
	for name := range ids {
		names = append(names, name)
	}

	f := func(tx *dataprovider.Tx) error {
		cityStore := NewCityStore(s.db, s.txer).WithTx(tx)
		cityFilter := dataprovider.NewCityFilter().ByNames(names...)
		cities, err := cityStore.ListByFilter(ctx, cityFilter)
		if err != nil {
			return errors.Wrap(err, "getting cities from city store")
		}
		if len(cities) == 0 {
			return errors.New("found no city")
		}
		for _, city := range cities {
			ids[city.Name] = city.ID
		}
		ib := sq.Insert("bus").Columns("city_id", "num")
		for _, bus := range buses {
			id := ids[bus.City]
			if id < 0 {
				s.logger(ctx).Debugf("bus [%s, %s] skipped", bus.City, bus.Num)
				continue
			}
			ib = ib.Values(id, bus.Num)
		}

		query, args, err := ib.ToSql()
		if err != nil {
			return errors.Wrap(err, "creating sql query for inserting buses")
		}

		s.logger(ctx).
			WithFields(
				"query", query,
				"args", args).
			Debug("inserting buses")

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return errors.Wrapf(err, "inserting buses with query %s", query)
		}
		return nil
	}
	return dataprovider.BeginAutoCommitedTx(ctx, s.txer, f)
}
