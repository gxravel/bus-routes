package handler

import (
	"encoding/json"
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
	"github.com/gxravel/bus-routes/internal/logger"
	v1 "github.com/gxravel/bus-routes/internal/model/v1"
	"github.com/pkg/errors"
)

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

	var cities = make([]*v1.City, 0)
	if err := json.NewDecoder(r.Body).Decode(&cities); err != nil {
		logger.FromContext(ctx).WithErr(err).Error("decoding data from post cities request")
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if len(cities) == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, errors.New("must provide cities"))
		return
	}

	err := s.busroutes.PostCities(ctx, cities...)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondCreated(w)
}

func (s *Server) putCity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var city = &v1.City{}
	if err := json.NewDecoder(r.Body).Decode(&city); err != nil {
		logger.FromContext(ctx).WithErr(err).Error("decoding data from put city request")
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if city.ID == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, errors.New("must provide city"))
		return
	}

	err := s.busroutes.PutCity(ctx, city)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondNoContent(w)
}

func (s *Server) deleteCity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter, err := api.ParseDeleteCityFilter(r)
	if err != nil || len(filter.IDs) == 0 && len(filter.Names) == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	err = s.busroutes.DeleteCity(ctx, filter)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondNoContent(w)
}
