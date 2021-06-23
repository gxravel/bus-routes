package busroutes

import (
	"context"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/model"
	v1 "github.com/gxravel/bus-routes/internal/model/v1"
)

func (r *BusRoutes) GetBuses(ctx context.Context, filter *dataprovider.BusFilter) ([]*v1.Bus, error) {
	dbBuses, err := r.busStore.ListByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return buses(dbBuses...), nil
}

func (r *BusRoutes) PostBuses(ctx context.Context, buses ...*v1.Bus) error {
	return r.busStore.New(ctx, dbBuses(buses...)...)
}

func dbBuses(buses ...*v1.Bus) []*model.Bus {
	var dbBuses = make([]*model.Bus, 0, len(buses))
	for _, bus := range buses {
		dbBuses = append(dbBuses, &model.Bus{
			City: bus.City,
			Num:  bus.Num,
		})
	}
	return dbBuses
}

func buses(dbBuses ...*model.Bus) []*v1.Bus {
	var cities = make([]*v1.Bus, 0, len(dbBuses))
	for _, bus := range dbBuses {
		cities = append(cities, &v1.Bus{
			ID:   bus.ID,
			City: bus.City,
			Num:  bus.Num,
		})
	}
	return cities
}
