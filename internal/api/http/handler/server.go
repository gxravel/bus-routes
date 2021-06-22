package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gxravel/bus-routes/assets"
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

	if cfg.API.ServeSwagger {
		registerSwagger(r)
	}

	r.Route("/internal", func(r chi.Router) {
		r.Get("/health", srv.getHealth)
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/cities", func(r chi.Router) {
				r.Get("/", srv.getCities)
				r.Post("/", srv.postCities)
			})
			r.Route("/buses", func(r chi.Router) {
				r.Get("/", srv.getBuses)
				r.Post("/", srv.postBuses)
			})
		})
	})

	srv.Handler = r

	return srv
}

func registerSwagger(r *chi.Mux) {
	r.HandleFunc("/internal/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/internal/swagger/", http.StatusFound)
	})

	swaggerHandler := http.StripPrefix("/internal/", http.FileServer(http.FS(assets.SwaggerFiles)))
	r.Get("/internal/swagger/*", swaggerHandler.ServeHTTP)
}
