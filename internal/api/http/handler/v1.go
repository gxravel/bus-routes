package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
)

func (s Server) getBuses(w http.ResponseWriter, r *http.Request) {
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
