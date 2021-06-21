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
		eq["id"] = f.IDs
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
			"city.name",
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
		if err == sql.ErrNoRows {
			return buses, nil
		}
		return nil, errors.Wrapf(err, "selecting buses by filter from database with query %s", query)
	}

	return buses, nil
}
