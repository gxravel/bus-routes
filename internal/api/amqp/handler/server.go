package handler

import (
	"context"

	"github.com/gxravel/bus-routes/internal/busroutes"
	"github.com/gxravel/bus-routes/internal/config"
	log "github.com/gxravel/bus-routes/internal/logger"
	"github.com/gxravel/bus-routes/pkg/rmq"
)

type Server struct {
	logger    log.Logger
	busroutes *busroutes.Busroutes
	handlers  []func(context.Context)
}

// NewServer creates new instance of the Server and subscribes for the amqp events,
// linking handlers with the required meta and amqp.Delivery.
func NewServer(
	cfg config.Config,
	publisher *rmq.Publisher,
	consumer *rmq.Consumer,
	busroutes *busroutes.Busroutes,
	logger log.Logger,
) (*Server, error) {

	srv := &Server{
		logger:    logger.WithModule("api:amqp"),
		busroutes: busroutes,
	}

	srv.handlers = make([]func(context.Context), 1)

	var err error
	srv.handlers[0], err = consumer.ListenRPCForDetailedRoutes(srv.getDetailedRoutes, publisher)
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
