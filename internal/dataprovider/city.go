package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes/internal/model"
)

type CityStore interface {
	WithTx(*Tx) CityStore
	ByFilter(ctx context.Context, filter *CityFilter) (*model.City, error)
	ListByFilter(ctx context.Context, filter *CityFilter) ([]*model.City, error)
	New(ctx context.Context, names ...string) error
}

type CityFilter struct {
	IDs   []int
	Names []string
}

func NewCityFilter() *CityFilter {
	return &CityFilter{}
}

func (f *CityFilter) ByIDs(ids ...int) *CityFilter {
	f.IDs = ids
	return f
}

func (f *CityFilter) ByNames(names ...string) *CityFilter {
	f.Names = names
	return f
}
