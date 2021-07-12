package busroutes

import (
	"context"

	httpv1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/model"
)

func (r *Busroutes) GetCities(ctx context.Context, filter *dataprovider.CityFilter) ([]*httpv1.City, error) {
	dbCities, err := r.cityStore.GetListByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	return toV1Cities(dbCities...), nil
}

func (r *Busroutes) AddCities(ctx context.Context, cities ...*httpv1.City) error {
	return r.cityStore.Add(ctx, toDBCities(cities...)...)
}

func (r *Busroutes) UpdateCity(ctx context.Context, city *httpv1.City) error {
	return r.cityStore.Update(ctx, toDBCities(city)[0])
}

func (r *Busroutes) DeleteCity(ctx context.Context, filter *dataprovider.CityFilter) error {
	return r.cityStore.Delete(ctx, filter)
}

func toDBCities(cities ...*httpv1.City) []*model.City {
	var dbCities = make([]*model.City, 0, len(cities))
	for _, city := range cities {
		dbCities = append(dbCities, &model.City{
			ID:   city.ID,
			Name: city.Name,
		})
	}

	return dbCities
}

func toV1Cities(dbCities ...*model.City) []*httpv1.City {
	var cities = make([]*httpv1.City, 0, len(dbCities))
	for _, city := range dbCities {
		cities = append(cities, &httpv1.City{
			ID:   city.ID,
			Name: city.Name,
		})
	}

	return cities
}
