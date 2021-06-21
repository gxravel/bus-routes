package busroutes

import (
	"context"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/model"
)

func (r *BusRoutes) GetBuses(ctx context.Context, busFilter *dataprovider.BusFilter) ([]*model.Bus, error) {
	dbBuses, err := r.busStore.ListByFilter(ctx, busFilter)
	if err != nil {
		return nil, err
	}
	return dbBuses, nil
}
