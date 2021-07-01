package busroutes

import (
	"context"

	v1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/model"
)

func (r *BusRoutes) GetStops(ctx context.Context, filter *dataprovider.StopFilter) ([]*v1.Stop, error) {
	dbStops, err := r.stopStore.GetListByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return toV1Stops(dbStops...), nil
}

func (r *BusRoutes) AddStops(ctx context.Context, stops ...*v1.Stop) error {
	return r.stopStore.Add(ctx, toDBStops(stops...)...)
}

func (r *BusRoutes) UpdateStops(ctx context.Context, stop *v1.Stop) error {
	return r.stopStore.Update(ctx, toDBStops(stop)[0])
}

func (r *BusRoutes) DeleteStop(ctx context.Context, filter *dataprovider.StopFilter) error {
	return r.stopStore.Delete(ctx, filter)
}

func toDBStops(stops ...*v1.Stop) []*model.Stop {
	var dbStops = make([]*model.Stop, 0, len(stops))
	for _, stop := range stops {
		dbStops = append(dbStops, &model.Stop{
			ID:      stop.ID,
			City:    stop.City,
			Address: stop.Address,
		})
	}
	return dbStops
}

func toV1Stops(dbStops ...*model.Stop) []*v1.Stop {
	var stops = make([]*v1.Stop, 0, len(dbStops))
	for _, stop := range dbStops {
		stops = append(stops, &v1.Stop{
			ID:      stop.ID,
			City:    stop.City,
			Address: stop.Address,
		})
	}
	return stops
}
