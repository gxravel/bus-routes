package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes/internal/model"
)

type CityStore interface {
	WithTx(*Tx) CityStore
	GetByFilter(ctx context.Context, filter *CityFilter) (*model.City, error)
	GetListByFilter(ctx context.Context, filter *CityFilter) ([]*model.City, error)
	Add(ctx context.Context, cities ...*model.City) error
	Delete(ctx context.Context, filter *CityFilter) error
	Update(ctx context.Context, city *model.City) error
}

type CityFilter struct {
	IDs   []int
	Names []string
}

func NewCityFilter() *CityFilter {
	return &CityFilter{}
}

// ByIDs filters by city.id.
func (f *CityFilter) ByIDs(ids ...int) *CityFilter {
	f.IDs = ids
	return f
}

// ByNames filters by city.name.
func (f *CityFilter) ByNames(names ...string) *CityFilter {
	f.Names = names
	return f
}
