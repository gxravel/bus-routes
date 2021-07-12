package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
	httpv1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	ierr "github.com/gxravel/bus-routes/internal/errors"
)

var (
	errMustProvideStop = ierr.NewReason(ierr.ErrMustProvide).WithMessage("stop")
)

func (s *Server) getStops(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stopFilter, err := api.ParseStopFilter(r)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	stops, err := s.busroutes.GetStops(ctx, stopFilter)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondDataOK(ctx, w, httpv1.RangeItemsResponse{
		Items: stops,
		Total: int64(len(stops)),
	})
}

func (s *Server) addStops(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var stops = make([]*httpv1.Stop, 0)
	if err := s.processRequest(r, &stops); err != nil {
		api.RespondError(ctx, w, err)
		return
	}
	if len(stops) == 0 {
		api.RespondError(ctx, w, errMustProvideStop)
		return
	}

	if err := s.busroutes.AddStops(ctx, stops...); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondCreated(w)
}

func (s *Server) updateStop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var stop = &httpv1.Stop{}
	if err := s.processRequest(r, stop); err != nil {
		api.RespondError(ctx, w, err)
		return
	}
	if stop.ID == 0 {
		api.RespondError(ctx, w, errMustProvideStop)
		return
	}

	if err := s.busroutes.UpdateStops(ctx, stop); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondNoContent(w)
}

func (s *Server) deleteStop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter, err := api.ParseDeleteStopFilter(r)
	if err != nil || len(filter.IDs) == 0 && len(filter.Addresses) == 0 {
		api.RespondError(ctx, w, err)
		return
	}

	if err = s.busroutes.DeleteStop(ctx, filter); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondNoContent(w)
}
