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

// ByFilter returns route depend on received filters.
func (s *RouteStore) ByFilter(ctx context.Context, filter *dataprovider.RouteFilter) (*model.Route, error) {
	routes, err := s.ListByFilter(ctx, filter)

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

// ListByFilter returns routes depend on received filters.
func (s *RouteStore) ListByFilter(ctx context.Context, filter *dataprovider.RouteFilter) ([]*model.Route, error) {
	qb := sq.
		Select(
			"bus_id",
			"stop_id",
			"step",
		).
		From(s.tableName).
		Where(routeCond(filter)).
		OrderBy("step")

	result, err := selectContext(ctx, qb, s.tableName, s.db, TypeRoute)
	if err != nil {
		return nil, err
	}
	return result.([]*model.Route), nil
}

// Add creates new routes.
func (s *RouteStore) Add(ctx context.Context, routes ...*model.Route) error {
	qb := sq.Insert(s.tableName).Columns("bus_id", "stop_id", "step")
	for _, route := range routes {
		qb = qb.Values(route.BusID, route.StopID, route.Step)
	}
	return execContext(ctx, qb, s.tableName, s.txer)
}

// Update updates route's stop_id.
func (s *RouteStore) Update(ctx context.Context, route *model.Route) error {
	qb := sq.Update(s.tableName).Set("stop_id", route.StopID).Where(sq.Eq{"bus_id": route.BusID, "step": route.Step})
	return execContext(ctx, qb, s.tableName, s.txer)
}

// Delete deletes route depend on received filter.
func (s *RouteStore) Delete(ctx context.Context, filter *dataprovider.RouteFilter) error {
	qb := sq.Delete(s.tableName).Where(routeCond(filter))
	return execContext(ctx, qb, s.tableName, s.txer)
}
