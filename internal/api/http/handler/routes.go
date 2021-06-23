package handler

import (
	"encoding/json"
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
	"github.com/gxravel/bus-routes/internal/logger"
	v1 "github.com/gxravel/bus-routes/internal/model/v1"
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

func (s *Server) postRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var routes = make([]*v1.Route, 0)
	if err := json.NewDecoder(r.Body).Decode(&routes); err != nil {
		logger.FromContext(ctx).WithErr(err).Error("decoding data from post routes request")
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if len(routes) == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, errors.New("must provide routes"))
		return
	}

	err := s.busroutes.PostRoutes(ctx, routes...)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondCreated(w)
}

func (s *Server) putRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var route = &v1.Route{}
	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		logger.FromContext(ctx).WithErr(err).Error("decoding data from put route request")
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if route.BusID == 0 || route.Step == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, errors.New("must provide route"))
		return
	}

	err := s.busroutes.PutRoute(ctx, route)
	if err != nil {
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

	err = s.busroutes.DeleteRoute(ctx, filter)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondNoContent(w)
}
