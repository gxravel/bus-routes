package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
)

func (s *Server) getHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := s.busroutes.IsHealthy(ctx); err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondEmpty(w)
}
