package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
	httpv1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	ierr "github.com/gxravel/bus-routes/internal/errors"
)

var (
	errMustProvideRoute = ierr.NewReason(ierr.ErrMustProvide).WithMessage("route")
)

func (s *Server) getRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	routeFilter, err := api.ParseRouteFilter(r)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	routes, err := s.busroutes.GetRoutes(ctx, routeFilter)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondDataOK(ctx, w, httpv1.RangeItemsResponse{
		Items: routes,
		Total: int64(len(routes)),
	})
}

// getDetailedRoutes returns the routes detailed view: city, address, number instead of ids.
func (s *Server) getDetailedRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	routeFilter, err := api.ParseRouteDetailedFilter(r)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	routes, err := s.busroutes.GetDetailedRoutes(ctx, routeFilter)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondDataOK(ctx, w, httpv1.RangeItemsResponse{
		Items: routes,
		Total: int64(len(routes)),
	})
}

func (s *Server) addRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var routes = make([]*httpv1.Route, 0)
	if err := s.processRequest(r, &routes); err != nil {
		api.RespondError(ctx, w, err)
		return
	}
	if len(routes) == 0 {
		api.RespondError(ctx, w, errMustProvideRoute)
		return
	}

	if err := s.busroutes.AddRoutes(ctx, routes...); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondCreated(w)
}

func (s *Server) updateRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var route = &httpv1.Route{}
	if err := s.processRequest(r, route); err != nil {
		api.RespondError(ctx, w, err)
		return
	}
	if route.BusID == 0 || route.Step == 0 {
		api.RespondError(ctx, w, errMustProvideRoute)
		return
	}

	if err := s.busroutes.UpdateRoute(ctx, route); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondNoContent(w)
}

func (s *Server) deleteRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter, err := api.ParseDeleteRouteFilter(r)
	if err != nil || len(filter.BusIDs) == 0 && len(filter.Steps) == 0 {
		api.RespondError(ctx, w, err)
		return
	}

	if err = s.busroutes.DeleteRoute(ctx, filter); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondNoContent(w)
}
