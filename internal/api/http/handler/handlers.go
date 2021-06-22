package handler

import (
	"encoding/json"
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
	"github.com/gxravel/bus-routes/internal/logger"
	"github.com/gxravel/bus-routes/internal/model"
)

func (srv *Server) getHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := srv.busroutes.IsHealthy(ctx)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondEmpty(w)
}

func (s *Server) getBuses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	busFilter, err := api.ParseBusFilter(r)
	if err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	buses, err := s.busroutes.GetBuses(ctx, busFilter)
	if err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	api.RespondDataOK(ctx, w, api.RangeItemsResponse{
		Items: buses,
		Total: int64(len(buses)),
	})
}

func (s *Server) postBuses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var buses = make([]*model.Bus, 0)
	if err := json.NewDecoder(r.Body).Decode(&buses); err != nil {
		logger.FromContext(ctx).WithErr(err).Error("decoding data from post buses request")
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	err := s.busroutes.PostBuses(ctx, buses...)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondCreated(w)
}

func (s *Server) getCities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cityFilter, err := api.ParseCityFilter(r)
	if err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	cities, err := s.busroutes.GetCities(ctx, cityFilter)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondDataOK(ctx, w, api.RangeItemsResponse{
		Items: cities,
		Total: int64(len(cities)),
	})
}

func (s *Server) postCities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var names = make([]string, 0)
	if err := json.NewDecoder(r.Body).Decode(&names); err != nil {
		logger.FromContext(ctx).WithErr(err).Error("decoding data from post cities request")
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	err := s.busroutes.PostCities(ctx, names...)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondCreated(w)
}
