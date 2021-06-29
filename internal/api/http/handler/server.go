package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gxravel/bus-routes/assets"
	mw "github.com/gxravel/bus-routes/internal/api/http/middleware"
	"github.com/gxravel/bus-routes/internal/busroutes"
	"github.com/gxravel/bus-routes/internal/config"
	"github.com/gxravel/bus-routes/internal/logger"

	"github.com/go-chi/chi"
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
			r.Route("/auth", func(r chi.Router) {
				r.Post("/signup", srv.signup)
				r.Post("/login", srv.login)
			})
			r.Route("/cities", func(r chi.Router) {
				r.Get("/", srv.getCities)
				r.Post("/", srv.addCities)
				r.Put("/", srv.updateCity)
				r.Delete("/", srv.deleteCity)
			})
			r.Route("/buses", func(r chi.Router) {
				r.Get("/", srv.getBuses)
				r.Post("/", srv.addBuses)
			})
			r.Route("/stops", func(r chi.Router) {
				r.Get("/", srv.getStops)
				r.Post("/", srv.addStops)
				r.Put("/", srv.updateStop)
				r.Delete("/", srv.deleteStop)
			})
			r.Route("/routes", func(r chi.Router) {
				r.Get("/", srv.getRoutes)
				r.Post("/", srv.addRoutes)
				r.Put("/", srv.updateRoute)
				r.Delete("/", srv.deleteRoute)
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

func (s *Server) processRequest(r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		s.logger.WithErr(err).Error("decoding data")
		return err
	}
	return nil
}
