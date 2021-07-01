package busroutes

import (
	"context"

	v1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/model"
)

func (r *BusRoutes) GetRoutes(ctx context.Context, filter *dataprovider.RouteFilter) ([]*v1.Route, error) {
	dbRoutes, err := r.routeStore.GetListByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return toV1Routes(dbRoutes...), nil
}

func (r *BusRoutes) GetDetailedRoutes(ctx context.Context, filter *dataprovider.RouteFilter) ([]*v1.RouteDetailed, error) {
	dbRoutes, err := r.routeStore.GetListByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return toV1RoutesDetailed(dbRoutes...), nil
}

func (r *BusRoutes) AddRoutes(ctx context.Context, routes ...*v1.Route) error {
	return r.routeStore.Add(ctx, toDBRoutes(routes...)...)
}

func (r *BusRoutes) UpdateRoute(ctx context.Context, route *v1.Route) error {
	return r.routeStore.Update(ctx, toDBRoutes(route)[0])
}

func (r *BusRoutes) DeleteRoute(ctx context.Context, filter *dataprovider.RouteFilter) error {
	return r.routeStore.Delete(ctx, filter)
}

func toDBRoutes(routes ...*v1.Route) []*model.Route {
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

func toV1Routes(dbRoutes ...*model.Route) []*v1.Route {
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

// toV1RoutesDetailed converts dbRoutes to v1.RouteDetailed.
// It expects dbRoutes to be ordered by bus_id.
func toV1RoutesDetailed(dbRoutes ...*model.Route) []*v1.RouteDetailed {
	var routes = make([]*v1.RouteDetailed, 0)
	var busID int64
	var lastRoute int = -1

	for _, route := range dbRoutes {
		if busID != route.BusID {
			busID = route.BusID
			routes = append(routes, &v1.RouteDetailed{
				City:   route.City,
				Bus:    route.Number,
				Points: make([]v1.RoutePoint, 0),
			})
			lastRoute++
		}
		routes[lastRoute].Points = append(routes[lastRoute].Points, v1.RoutePoint{Step: route.Step, Address: route.Address})
	}
	return routes
}
