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

func (s *CityStore) New(ctx context.Context, cities ...*model.City) error {
	qb := sq.Insert(s.tableName).Columns("name")
	for _, city := range cities {
		qb = qb.Values(city.Name)
	}

	query, args, err := qb.ToSql()
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
			return errors.Wrapf(err, "inserting cities with query %s", query)
		}
		return nil
	}

	return dataprovider.BeginAutoCommitedTx(ctx, s.txer, f)
}

func (s *CityStore) Update(ctx context.Context, city *model.City) error {
	qb := sq.Update(s.tableName).Set("name", city.Name).Where(sq.Eq{"id": city.ID})

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "creating sql query for updating city")
	}

	s.logger(ctx).
		WithFields(
			"query", query,
			"args", args).
		Debug("updating city")

	f := func(tx *dataprovider.Tx) error {
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return errors.Wrapf(err, "updating city with query %s", query)
		}
		return nil
	}

	return dataprovider.BeginAutoCommitedTx(ctx, s.txer, f)
}

func (s *CityStore) Delete(ctx context.Context, filter *dataprovider.CityFilter) error {
	qb := sq.Delete(s.tableName).Where(cityCond(filter))

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "creating sql query for deleting city")
	}

	s.logger(ctx).
		WithFields(
			"query", query,
			"args", args).
		Debug("deleting city")

	f := func(tx *dataprovider.Tx) error {
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return errors.Wrapf(err, "deleting city with query %s", query)
		}
		return nil
	}

	return dataprovider.BeginAutoCommitedTx(ctx, s.txer, f)
}
