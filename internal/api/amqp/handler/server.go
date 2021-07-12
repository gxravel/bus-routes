package handler

import (
	"context"

	"github.com/gxravel/bus-routes/internal/busroutes"
	"github.com/gxravel/bus-routes/internal/config"
	log "github.com/gxravel/bus-routes/internal/logger"
	"github.com/gxravel/bus-routes/pkg/rmq"
)

type Server struct {
	broker    *rmq.Client
	logger    log.Logger
	busroutes *busroutes.BusRoutes
	handlers  []func(context.Context)
}

// NewServer creates new instance of the Server and subscribes for the amqp events,
// linking handlers with the required meta and amqp.Delivery.
func NewServer(
	cfg config.Config,
	busroutes *busroutes.BusRoutes,
	logger log.Logger,
) (*Server, error) {

	broker, err := rmq.NewClient(cfg.RabbitMQ.URL, logger)
	if err != nil {
		logger.WithErr(err).Fatal("failed to create RabbitMQ client")
		return nil, err
	}

	srv := &Server{
		broker:    broker,
		logger:    logger.WithModule("api:amqp"),
		busroutes: busroutes,
	}

	srv.handlers = make([]func(context.Context), 1)

	srv.handlers[0], err = broker.ListenRPCForDetailedRoutes(srv.getDetailedRoutes)
	if err != nil {
		return nil, err
	}

	return srv, nil
}

// ListenAndServe runs handlers to listen to message deliveries.
func (s *Server) ListenAndServe() {
	ctx := context.Background()

	for _, handler := range s.handlers {
		go handler(ctx)
	}
}

// CloseConnection closes connection of amqp client.
func (s *Server) CloseConnection() error {
	return s.broker.Close()
}
