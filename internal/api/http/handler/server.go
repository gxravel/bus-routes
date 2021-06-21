package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	mw "github.com/gxravel/bus-routes/internal/api/http/middleware"
	"github.com/gxravel/bus-routes/internal/busroutes"
	"github.com/gxravel/bus-routes/internal/config"
	"github.com/gxravel/bus-routes/internal/logger"
)

type Server struct {
	*http.Server
	logger    logger.Logger
	busroutes *busroutes.BusRoutes
}

func NewServer(
	cfg *config.Config,
	busroutes *busroutes.BusRoutes,
	logger logger.Logger,
) *Server {
	srv := &Server{
		Server: &http.Server{
			Addr:         cfg.API.Address,
			ReadTimeout:  cfg.API.ReadTimeout,
			WriteTimeout: cfg.API.WriteTimeout,
		},
		logger:    logger.WithStr("module", "api:http"),
		busroutes: busroutes,
	}

	r := chi.NewRouter()

	r.Use(mw.Logger(srv.logger))
	r.Use(mw.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/buses", func(r chi.Router) {
				r.Get("/", srv.getBuses)
			})
		})
	})

	srv.Handler = r

	return srv
}
