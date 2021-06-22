package busroutes

import (
	"context"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/model"
)

func (r *BusRoutes) IsHealthy(ctx context.Context) error {
	return r.db.StatusCheck(ctx)
}

func (r *BusRoutes) GetBuses(ctx context.Context, busFilter *dataprovider.BusFilter) ([]*model.Bus, error) {
	dbBuses, err := r.busStore.ListByFilter(ctx, busFilter)
	if err != nil {
		return nil, err
	}
	return dbBuses, nil
}
func (r *BusRoutes) PostBuses(ctx context.Context, buses ...*model.Bus) error {
	return r.busStore.New(ctx, buses...)
}

func (r *BusRoutes) GetCities(ctx context.Context, cityFilter *dataprovider.CityFilter) ([]*model.City, error) {
	return r.cityStore.ListByFilter(ctx, cityFilter)
}
func (r *BusRoutes) PostCities(ctx context.Context, names ...string) error {
	return r.cityStore.New(ctx, names...)
}
