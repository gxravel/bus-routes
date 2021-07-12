package handler

import (
	"context"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/pkg/rmq"
	v1 "github.com/gxravel/bus-routes/pkg/rmq/v1"

	"github.com/streadway/amqp"
)

// getDetailedRoutes returns the routes detailed view: city, address, number.
func (s *Server) getDetailedRoutes(ctx context.Context, message *amqp.Delivery) (interface{}, error) {
	s.logger.WithField("delivery", message).Debug("got message")

	bus := &v1.Bus{}
	if err := rmq.TranslateMessage(message.Body, bus); err != nil {
		return nil, err
	}

	buses, err := s.busroutes.GetBuses(ctx, dataprovider.NewBusFilter().
		ByCities(bus.City).
		ByNums(bus.Num),
	)
	if err != nil {
		return nil, err
	}
	if len(buses) == 0 {
		return nil, nil
	}

	routesDetailed, err := s.busroutes.GetDetailedRoutes(ctx, dataprovider.NewRouteFilter().
		ByBusIDs(buses[0].ID).
		ViewDetailed(),
	)
	if err != nil {
		return nil, err
	}
	if len(routesDetailed) == 0 {
		return nil, nil
	}

	return v1.RangeItemsResponse{
		Items: routesDetailed,
		Total: int64(len(routesDetailed)),
	}, nil
}
