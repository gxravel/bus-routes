package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
	httpv1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	ierr "github.com/gxravel/bus-routes/internal/errors"
)

var (
	errMustProvideCity = ierr.NewReason(ierr.ErrMustProvide).WithMessage("city")
)

func (s *Server) getCities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cityFilter, err := api.ParseCityFilter(r)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	cities, err := s.busroutes.GetCities(ctx, cityFilter)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondDataOK(ctx, w, httpv1.RangeItemsResponse{
		Items: cities,
		Total: int64(len(cities)),
	})
}

func (s *Server) addCities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var cities = make([]*httpv1.City, 0)
	if err := s.processRequest(r, &cities); err != nil {
		api.RespondError(ctx, w, err)
		return
	}
	if len(cities) == 0 {
		api.RespondError(ctx, w, errMustProvideCity)
		return
	}

	if err := s.busroutes.AddCities(ctx, cities...); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondCreated(w)
}

func (s *Server) updateCity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var city = &httpv1.City{}
	if err := s.processRequest(r, city); err != nil {
		api.RespondError(ctx, w, err)
		return
	}
	if city.ID == 0 {
		api.RespondError(ctx, w, errMustProvideCity)
		return
	}

	if err := s.busroutes.UpdateCity(ctx, city); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondNoContent(w)
}

func (s *Server) deleteCity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter, err := api.ParseDeleteCityFilter(r)
	if err != nil || len(filter.IDs) == 0 && len(filter.Names) == 0 {
		api.RespondError(ctx, w, err)
		return
	}

	if err = s.busroutes.DeleteCity(ctx, filter); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondNoContent(w)
}
