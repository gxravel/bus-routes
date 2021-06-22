package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes/internal/model"
)

type BusStore interface {
	WithTx(*Tx) BusStore
	ByFilter(ctx context.Context, filter *BusFilter) (*model.Bus, error)
	ListByFilter(ctx context.Context, filter *BusFilter) ([]*model.Bus, error)
	New(ctx context.Context, buses ...*model.Bus) error
}

type BusFilter struct {
	IDs       []int64
	Cities    []string
	Nums      []string
	Paginator *Paginator
}

func NewBusFilter() *BusFilter {
	return &BusFilter{}
}

func (f *BusFilter) ByIDs(ids ...int64) *BusFilter {
	f.IDs = ids
	return f
}

func (f *BusFilter) ByCities(cities ...string) *BusFilter {
	f.Cities = cities
	return f
}

func (f *BusFilter) ByNums(nums ...string) *BusFilter {
	f.Nums = nums
	return f
}

// WithPaginator adds pagination to filter.
func (f *BusFilter) WithPaginator(paginator *Paginator) *BusFilter {
	f.Paginator = paginator
	return f
}
