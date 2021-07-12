package mysql

import (
	"context"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/model"
	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// RouteStore is route mysql store.
type RouteStore struct {
	db        sqlx.ExtContext
	txer      dataprovider.Txer
	tableName string
}

// NewRouteStore creates new instance of RouteStore.
func NewRouteStore(db sqlx.ExtContext, txer dataprovider.Txer) *RouteStore {
	return &RouteStore{
		db:        db,
		txer:      txer,
		tableName: "route",
	}
}

// WithTx sets transaction as active connection.
func (s *RouteStore) WithTx(tx *dataprovider.Tx) dataprovider.RouteStore {
	return &RouteStore{
		db:        tx,
		tableName: s.tableName,
	}
}

func routeCond(f *dataprovider.RouteFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	var cond sq.Sqlizer = eq

	if len(f.BusIDs) > 0 {
		eq["route.bus_id"] = f.BusIDs
	}
	if len(f.StopIDs) > 0 {
		eq["route.stop_id"] = f.StopIDs
	}
	if len(f.Steps) > 0 {
		eq["route.step"] = f.Steps
	}

	return cond
}

func (s *RouteStore) columns(filter *dataprovider.RouteFilter) []string {
	if filter != nil && filter.DetailedView {
		return []string{
			"route.bus_id as bus_id",
			"city.name as city",
			"num",
			"step",
			"address",
		}
	}
	return []string{
		"bus_id",
		"stop_id",
		"step",
	}
}

func (s *RouteStore) joins(qb sq.SelectBuilder, filter *dataprovider.RouteFilter) sq.SelectBuilder {
	if filter.DetailedView {
		qb = qb.Join("bus ON route.bus_id = bus.id").
			Join("stop ON route.stop_id = stop.id").
			Join("city ON bus.city_id = city.id")
	}
	return qb
}

func (s *RouteStore) ordersBy(qb sq.SelectBuilder, filter *dataprovider.RouteFilter) sq.SelectBuilder {
	qb = qb.OrderBy("bus_id", "step")
	return qb
}

// GetByFilter returns route depend on received filters.
func (s *RouteStore) GetByFilter(ctx context.Context, filter *dataprovider.RouteFilter) (*model.Route, error) {
	routes, err := s.GetListByFilter(ctx, filter)

	switch {
	case err != nil:
		return nil, err
	case len(routes) == 0:
		return nil, nil
	case len(routes) == 1:
		return routes[0], nil
	default:
		return nil, errors.New("fetched more than 1 route")
	}
}

// GetListByFilter returns routes depend on received filters.
func (s *RouteStore) GetListByFilter(ctx context.Context, filter *dataprovider.RouteFilter) ([]*model.Route, error) {
	qb := sq.
		Select(s.columns(filter)...).
		From(s.tableName).
		Where(routeCond(filter))

	qb = s.joins(qb, filter)
	qb = s.ordersBy(qb, filter)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	message := "select " + s.tableName + " by filter with query " + query

	var result = make([]*model.Route, 0)
	if err := sqlx.SelectContext(ctx, s.db, &result, query, args...); err != nil {
		return nil, errors.Wrapf(err, message)
	}

	return result, nil
}

// Add creates new routes.
func (s *RouteStore) Add(ctx context.Context, routes ...*model.Route) error {
	qb := sq.Insert(s.tableName).Columns(s.columns(nil)...)
	for _, route := range routes {
		qb = qb.Values(route.BusID, route.StopID, route.Step)
	}

	return execContext(ctx, qb, s.tableName, s.db)
}

// Update updates route's stop_id.
func (s *RouteStore) Update(ctx context.Context, route *model.Route) error {
	qb := sq.Update(s.tableName).
		Set("stop_id", route.StopID).
		Where(sq.Eq{
			"bus_id": route.BusID,
			"step":   route.Step},
		)

	return execContext(ctx, qb, s.tableName, s.db)
}

// Delete deletes route depend on received filter.
func (s *RouteStore) Delete(ctx context.Context, filter *dataprovider.RouteFilter) error {
	qb := sq.Delete(s.tableName).Where(routeCond(filter))
	return execContext(ctx, qb, s.tableName, s.db)
}
