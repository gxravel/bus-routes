package busroutes

import (
	"context"

	v1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/model"
)

func (r *BusRoutes) GetRoutes(ctx context.Context, filter *dataprovider.RouteFilter) ([]*v1.Route, error) {
	dbRoutes, err := r.routeStore.ListByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return routes(dbRoutes...), nil
}

func (r *BusRoutes) AddRoutes(ctx context.Context, routes ...*v1.Route) error {
	return r.routeStore.Add(ctx, dbRoutes(routes...)...)
}

func (r *BusRoutes) UpdateRoute(ctx context.Context, route *v1.Route) error {
	return r.routeStore.Update(ctx, dbRoutes(route)[0])
}

func (r *BusRoutes) DeleteRoute(ctx context.Context, filter *dataprovider.RouteFilter) error {
	return r.routeStore.Delete(ctx, filter)
}

func dbRoutes(routes ...*v1.Route) []*model.Route {
	var dbRoutes = make([]*model.Route, 0, len(routes))
	for _, route := range routes {
		dbRoutes = append(dbRoutes, &model.Route{
			BusID:  route.BusID,
			StopID: route.StopID,
			Step:   route.Step,
		})
	}
	return dbRoutes
}

func routes(dbRoutes ...*model.Route) []*v1.Route {
	var routes = make([]*v1.Route, 0, len(dbRoutes))
	for _, route := range dbRoutes {
		routes = append(routes, &v1.Route{
			BusID:  route.BusID,
			StopID: route.StopID,
			Step:   route.Step,
		})
	}
	return routes
}
