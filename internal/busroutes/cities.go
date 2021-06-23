package busroutes

import (
	"context"

	v1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/model"
)

func (r *BusRoutes) GetCities(ctx context.Context, filter *dataprovider.CityFilter) ([]*v1.City, error) {
	dbCities, err := r.cityStore.ListByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return cities(dbCities...), nil
}

func (r *BusRoutes) AddCities(ctx context.Context, cities ...*v1.City) error {
	return r.cityStore.Add(ctx, dbCities(cities...)...)
}

func (r *BusRoutes) UpdateCity(ctx context.Context, city *v1.City) error {
	return r.cityStore.Update(ctx, dbCities(city)[0])
}

func (r *BusRoutes) DeleteCity(ctx context.Context, filter *dataprovider.CityFilter) error {
	return r.cityStore.Delete(ctx, filter)
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

func cities(dbCities ...*model.City) []*v1.City {
	var cities = make([]*v1.City, 0, len(dbCities))
	for _, city := range dbCities {
		cities = append(cities, &v1.City{
			ID:   city.ID,
			Name: city.Name,
		})
	}
	return cities
}
