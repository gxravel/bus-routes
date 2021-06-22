package busroutes

import (
	"context"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/model"
	v1 "github.com/gxravel/bus-routes/internal/model/v1"
)

func (r *BusRoutes) IsHealthy(ctx context.Context) error {
	return r.db.StatusCheck(ctx)
}

func dbBuses(buses []*v1.Bus) []*model.Bus {
	var dbBuses = make([]*model.Bus, 0, len(buses))
	for _, bus := range buses {
		dbBuses = append(dbBuses, &model.Bus{
			City: bus.City,
			Num:  bus.Num,
		})
	}
	return dbBuses
}

func buses(dbBuses []*model.Bus) []*v1.Bus {
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

func (r *BusRoutes) GetBuses(ctx context.Context, filter *dataprovider.BusFilter) ([]*v1.Bus, error) {
	dbBuses, err := r.busStore.ListByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return buses(dbBuses), nil
}
func (r *BusRoutes) PostBuses(ctx context.Context, buses ...*v1.Bus) error {
	return r.busStore.New(ctx, dbBuses(buses)...)
}

func dbCities(cities ...*v1.City) []*model.City {
	var dbCities = make([]*model.City, 0, len(cities))
	for _, city := range cities {
		dbCities = append(dbCities, &model.City{
			ID:   city.ID,
			Name: city.Name,
		})
	}
	return dbCities
}

func cities(dbCities []*model.City) []*v1.City {
	var cities = make([]*v1.City, 0, len(dbCities))
	for _, city := range dbCities {
		cities = append(cities, &v1.City{
			ID:   city.ID,
			Name: city.Name,
		})
	}
	return cities
}

func (r *BusRoutes) GetCities(ctx context.Context, filter *dataprovider.CityFilter) ([]*v1.City, error) {
	dbCities, err := r.cityStore.ListByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return cities(dbCities), nil
}
func (r *BusRoutes) PostCities(ctx context.Context, cities ...*v1.City) error {
	err := r.cityStore.New(ctx, dbCities(cities...)...)
	return err
}
func (r *BusRoutes) PutCity(ctx context.Context, city *v1.City) error {
	err := r.cityStore.Update(ctx, dbCities(city)[0])
	return err
}
func (r *BusRoutes) DeleteCity(ctx context.Context, filter *dataprovider.CityFilter) error {
	return r.cityStore.Delete(ctx, filter)
}
