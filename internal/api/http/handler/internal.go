package handler

import (
	"net/http"

	api "github.com/gxravel/bus-routes/internal/api/http"
)

func (srv *Server) getHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := srv.busroutes.IsHealthy(ctx)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondEmpty(w)
}
