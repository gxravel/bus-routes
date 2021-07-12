package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
	httpv1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	ierr "github.com/gxravel/bus-routes/internal/errors"
)

func (s *Server) getBuses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	busFilter, err := api.ParseBusFilter(r)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	buses, err := s.busroutes.GetBuses(ctx, busFilter)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondDataOK(ctx, w, httpv1.RangeItemsResponse{
		Items: buses,
		Total: int64(len(buses)),
	})
}

func (s *Server) addBuses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var buses = make([]*httpv1.Bus, 0)
	if err := s.processRequest(r, &buses); err != nil {
		api.RespondError(ctx, w, err)
		return
	}
	if len(buses) == 0 {
		api.RespondError(ctx, w, ierr.NewReason(ierr.ErrMustProvide).WithMessage("buses"))
		return
	}

	if err := s.busroutes.AddBuses(ctx, buses...); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.RespondCreated(w)
}
