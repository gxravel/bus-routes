package handler

import (
	"encoding/json"
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
)

func (s *Server) getHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := s.busroutes.IsHealthy(ctx)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondEmpty(w)
}

func (s *Server) processRequest(r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		s.logger.WithErr(err).Error("decoding data")
		return err
	}
	return nil
}
