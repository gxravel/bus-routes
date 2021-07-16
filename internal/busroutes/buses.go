package busroutes

import (
	"context"

	httpv1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/model"
)

func (r *Busroutes) GetBuses(ctx context.Context, filter *dataprovider.BusFilter) ([]*httpv1.Bus, error) {
	dbBuses, err := r.busStore.GetListByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	return toV1Buses(dbBuses...), nil
}

func (r *Busroutes) AddBuses(ctx context.Context, buses ...*httpv1.Bus) error {
	return r.busStore.Add(ctx, toDBBuses(buses...)...)
}

func toDBBuses(buses ...*httpv1.Bus) []*model.Bus {
	var dbBuses = make([]*model.Bus, 0, len(buses))
	for _, bus := range buses {
		dbBuses = append(dbBuses, &model.Bus{
			City:   bus.City,
			Number: bus.Num,
		})
	}

	return dbBuses
}

func toV1Buses(dbBuses ...*model.Bus) []*httpv1.Bus {
	var buses = make([]*httpv1.Bus, 0, len(dbBuses))
	for _, bus := range dbBuses {
		buses = append(buses, &httpv1.Bus{
			ID:   bus.ID,
			City: bus.City,
			Num:  bus.Number,
		})
	}

	return buses
}
