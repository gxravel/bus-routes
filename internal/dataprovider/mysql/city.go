package mysql

import (
	"context"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// CityStore is city mysql store.
type CityStore struct {
	db        sqlx.ExtContext
	txer      dataprovider.Txer
	tableName string
}

// NewCityStore creates new instance of CityStore.
func NewCityStore(db sqlx.ExtContext, txer dataprovider.Txer) *CityStore {
	return &CityStore{
		db:        db,
		txer:      txer,
		tableName: "city",
	}
}

// WithTx sets transaction as active connection.
func (s *CityStore) WithTx(tx *dataprovider.Tx) dataprovider.CityStore {
	return &CityStore{
		db:        tx,
		tableName: s.tableName,
	}
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

// GetByFilter returns city depend on received filters.
func (s *CityStore) GetByFilter(ctx context.Context, filter *dataprovider.CityFilter) (*model.City, error) {
	cities, err := s.GetListByFilter(ctx, filter)

	switch {
	case err != nil:
		return nil, err
	case len(cities) == 0:
		return nil, nil
	case len(cities) == 1:
		return cities[0], nil
	default:
		return nil, errors.New("fetched more than 1 city")
	}
}

// GetListByFilter returns cities depend on received filters.
func (s *CityStore) GetListByFilter(ctx context.Context, filter *dataprovider.CityFilter) ([]*model.City, error) {
	qb := sq.
		Select(
			"id",
			"name",
		).
		From(s.tableName).
		Where(cityCond(filter))

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	message := "select " + s.tableName + " by filter with query " + query

	var result = make([]*model.City, 0)
	if err := sqlx.SelectContext(ctx, s.db, &result, query, args...); err != nil {
		return nil, errors.Wrapf(err, message)
	}

	return result, nil
}

// Add creates new cities.
func (s *CityStore) Add(ctx context.Context, cities ...*model.City) error {
	qb := sq.Insert(s.tableName).Columns("name")
	for _, city := range cities {
		qb = qb.Values(city.Name)
	}

	return execContext(ctx, qb, s.tableName, s.db)
}

// Update updates city name.
func (s *CityStore) Update(ctx context.Context, city *model.City) error {
	qb := sq.Update(s.tableName).
		Set("name", city.Name).
		Where(sq.Eq{"id": city.ID})

	return execContext(ctx, qb, s.tableName, s.db)
}

// Delete deletes city depend on received filter.
func (s *CityStore) Delete(ctx context.Context, filter *dataprovider.CityFilter) error {
	qb := sq.Delete(s.tableName).Where(cityCond(filter))
	return execContext(ctx, qb, s.tableName, s.db)
}

// getCitiesIDs return the ids as a map of names.
func getCitiesIDs(ctx context.Context, ids map[string]int, db sqlx.ExtContext, txer dataprovider.Txer, tx *dataprovider.Tx) error {
	var names = make([]string, 0, len(ids))
	for name := range ids {
		names = append(names, name)
	}

	cityStore := NewCityStore(db, txer).WithTx(tx)
	cityFilter := dataprovider.NewCityFilter().ByNames(names...)

	cities, err := cityStore.GetListByFilter(ctx, cityFilter)
	if err != nil {
		return errors.Wrap(err, "getting cities from city store")
	}
	if len(cities) == 0 {
		return errors.New("did not find a city")
	}

	for _, city := range cities {
		ids[city.Name] = city.ID
	}

	return nil
}

// getCityID returns the id by name.
func getCityID(ctx context.Context, name string, db sqlx.ExtContext, txer dataprovider.Txer, tx *dataprovider.Tx) (int, error) {
	cityStore := NewCityStore(db, txer).WithTx(tx)
	cityFilter := dataprovider.NewCityFilter().ByNames(name)

	city, err := cityStore.GetByFilter(ctx, cityFilter)
	if err != nil {
		return 0, errors.Wrap(err, "getting city from city store")
	}

	return city.ID, nil
}
