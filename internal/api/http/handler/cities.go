package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
	v1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"

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

func (s *Server) addCities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var cities = make([]*v1.City, 0)
	if err := s.processRequest(r, &cities); err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if len(cities) == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, errors.New("must provide cities"))
		return
	}

	err := s.busroutes.AddCities(ctx, cities...)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondCreated(w)
}

func (s *Server) updateCity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var city = &v1.City{}
	if err := s.processRequest(r, city); err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if city.ID == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, errors.New("must provide city"))
		return
	}

	err := s.busroutes.UpdateCity(ctx, city)
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
