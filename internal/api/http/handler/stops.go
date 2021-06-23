package handler

import (
	"encoding/json"
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
	"github.com/gxravel/bus-routes/internal/logger"
	v1 "github.com/gxravel/bus-routes/internal/model/v1"
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

func (s *Server) postStops(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var stops = make([]*v1.Stop, 0)
	if err := json.NewDecoder(r.Body).Decode(&stops); err != nil {
		logger.FromContext(ctx).WithErr(err).Error("decoding data from post stops request")
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if len(stops) == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, errors.New("must provide stops"))
		return
	}

	err := s.busroutes.PostStops(ctx, stops...)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondCreated(w)
}

func (s *Server) putStop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var stop = &v1.Stop{}
	if err := json.NewDecoder(r.Body).Decode(&stop); err != nil {
		logger.FromContext(ctx).WithErr(err).Error("decoding data from put stop request")
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if stop.ID == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, errors.New("must provide stop"))
		return
	}

	err := s.busroutes.PutStop(ctx, stop)
	if err != nil {
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

	err = s.busroutes.DeleteStop(ctx, filter)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondNoContent(w)
}
