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

type CityStore struct {
	db        sqlx.ExtContext
	txer      dataprovider.Txer
	tableName string
}

func NewCityStore(db sqlx.ExtContext, txer dataprovider.Txer) *CityStore {
	return &CityStore{
		db:        db,
		txer:      txer,
		tableName: "city",
	}
}

func (s *CityStore) WithTx(tx *dataprovider.Tx) dataprovider.CityStore {
	return &CityStore{
		db:        tx,
		tableName: s.tableName,
	}
}

func (s *CityStore) logger(ctx context.Context) logger.Logger {
	return logger.FromContext(ctx).WithStr("module", "mysql:city")
}

func cityCond(f *dataprovider.CityFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	var cond sq.Sqlizer = eq

	if len(f.IDs) > 0 {
		eq["city.id"] = f.IDs
	}
	if len(f.Names) > 0 {
		eq["name"] = f.Names
	}

	return cond
}

func (s *CityStore) ByFilter(ctx context.Context, filter *dataprovider.CityFilter) (*model.City, error) {
	cities, err := s.ListByFilter(ctx, filter)

	switch {
	case err != nil:
		return nil, err
	case len(cities) == 0:
		return nil, nil
	case len(cities) == 1:
		return cities[0], nil
	default:
		return nil, errors.New("fetched more than 1 City")
	}
}

func (s *CityStore) ListByFilter(ctx context.Context, filter *dataprovider.CityFilter) ([]*model.City, error) {
	qb := sq.
		Select(
			"id",
			"name",
		).
		From(s.tableName).
		Where(cityCond(filter))

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "creating sql query for getting cities by filter")
	}

	s.logger(ctx).
		WithFields(
			"query", query,
			"args", args).
		Debug("selecting cities by filter query SQL")

	cities := make([]*model.City, 0)
	if err = sqlx.SelectContext(ctx, s.db, &cities, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return cities, nil
		}
		return nil, errors.Wrapf(err, "selecting cities by filter from database with query %s", query)
	}

	return cities, nil
}

func (s *CityStore) New(ctx context.Context, names ...string) error {
	if len(names) == 0 {
		return nil
	}

	ib := sq.Insert("city").Columns("name")
	for _, name := range names {
		ib = ib.Values(name)
	}

	query, args, err := ib.ToSql()
	if err != nil {
		return errors.Wrap(err, "creating sql query for inserting cities")
	}

	s.logger(ctx).
		WithFields(
			"query", query,
			"args", args).
		Debug("inserting cities")
	f := func(tx *dataprovider.Tx) error {
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return errors.Wrapf(err, "inserting cities with query %s", query)
		}
		return nil
	}

	return dataprovider.BeginAutoCommitedTx(ctx, s.txer, f)
}
