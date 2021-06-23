package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
	v1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"

	"github.com/pkg/errors"
)

func (s *Server) getRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	routeFilter, err := api.ParseRouteFilter(r)
	if err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	routes, err := s.busroutes.GetRoutes(ctx, routeFilter)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondDataOK(ctx, w, api.RangeItemsResponse{
		Items: routes,
		Total: int64(len(routes)),
	})
}

func (s *Server) addRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var routes = make([]*v1.Route, 0)
	if err := s.processRequest(r, &routes); err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if len(routes) == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, errors.New("must provide routes"))
		return
	}

	if err := s.busroutes.AddRoutes(ctx, routes...); err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondCreated(w)
}

func (s *Server) updateRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var route = &v1.Route{}
	if err := s.processRequest(r, route); err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if route.BusID == 0 || route.Step == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, errors.New("must provide route"))
		return
	}

	if err := s.busroutes.UpdateRoute(ctx, route); err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondNoContent(w)
}

func (s *Server) deleteRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter, err := api.ParseDeleteRouteFilter(r)
	if err != nil || len(filter.BusIDs) == 0 && len(filter.Steps) == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if err = s.busroutes.DeleteRoute(ctx, filter); err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondNoContent(w)
}
