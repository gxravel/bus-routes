package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
	v1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"

	"github.com/pkg/errors"
)

func (s *Server) getStops(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stopFilter, err := api.ParseStopFilter(r)
	if err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	stops, err := s.busroutes.GetStops(ctx, stopFilter)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondDataOK(ctx, w, api.RangeItemsResponse{
		Items: stops,
		Total: int64(len(stops)),
	})
}

func (s *Server) addStops(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var stops = make([]*v1.Stop, 0)
	if err := s.processRequest(r, &stops); err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if len(stops) == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, errors.New("must provide stops"))
		return
	}

	if err := s.busroutes.AddStops(ctx, stops...); err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondCreated(w)
}

func (s *Server) updateStop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var stop = &v1.Stop{}
	if err := s.processRequest(r, stop); err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if stop.ID == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, errors.New("must provide stop"))
		return
	}

	if err := s.busroutes.UpdateStops(ctx, stop); err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondNoContent(w)
}

func (s *Server) deleteStop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter, err := api.ParseDeleteStopFilter(r)
	if err != nil || len(filter.IDs) == 0 && len(filter.Addresses) == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if err = s.busroutes.DeleteStop(ctx, filter); err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondNoContent(w)
}
