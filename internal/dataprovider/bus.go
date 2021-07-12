package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes/internal/model"
)

type BusStore interface {
	WithTx(*Tx) BusStore
	GetByFilter(ctx context.Context, filter *BusFilter) (*model.Bus, error)
	GetListByFilter(ctx context.Context, filter *BusFilter) ([]*model.Bus, error)
	Add(ctx context.Context, buses ...*model.Bus) error
}

type BusFilter struct {
	IDs         []int64
	Cities      []string
	CitiesIDs   []int
	Nums        []string
	DoPreferIDs bool
	Paginator   *Paginator
}

func NewBusFilter() *BusFilter {
	return &BusFilter{}
}

// ByIDs filters by bus.id.
func (f *BusFilter) ByIDs(ids ...int64) *BusFilter {
	f.IDs = ids
	return f
}

// ByCities filters by city.name.
func (f *BusFilter) ByCities(cities ...string) *BusFilter {
	f.Cities = cities
	return f
}

// ByCitiesIDs filters by bus.city_id.
func (f *BusFilter) ByCitiesIDs(citiesIDs ...int) *BusFilter {
	f.CitiesIDs = citiesIDs
	return f
}

// ByNums filters bu bus.num.
func (f *BusFilter) ByNums(nums ...string) *BusFilter {
	f.Nums = nums
	return f
}

// WithPaginator adds pagination to filter.
func (f *BusFilter) WithPaginator(paginator *Paginator) *BusFilter {
	f.Paginator = paginator
	return f
}

// PreferIDs select ids instead of joined values.
func (f *BusFilter) PreferIDs() *BusFilter {
	f.DoPreferIDs = true
	return f
}
