package handler

import (
	"encoding/json"
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
	"github.com/gxravel/bus-routes/internal/logger"
	v1 "github.com/gxravel/bus-routes/internal/model/v1"
	"github.com/pkg/errors"
)

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

	var buses = make([]*v1.Bus, 0)
	if err := json.NewDecoder(r.Body).Decode(&buses); err != nil {
		logger.FromContext(ctx).WithErr(err).Error("decoding data from post buses request")
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if len(buses) == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, errors.New("must provide buses"))
		return
	}

	err := s.busroutes.PostBuses(ctx, buses...)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondCreated(w)
}
