package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes/internal/model"
)

type StopStore interface {
	WithTx(*Tx) StopStore
	ByFilter(ctx context.Context, filter *StopFilter) (*model.Stop, error)
	ListByFilter(ctx context.Context, filter *StopFilter) ([]*model.Stop, error)
	New(ctx context.Context, cities ...*model.Stop) error
	Update(ctx context.Context, city *model.Stop) error
	Delete(ctx context.Context, filter *StopFilter) error
}

type StopFilter struct {
	IDs       []int64
	Cities    []string
	Addresses []string
}

func NewStopFilter() *StopFilter {
	return &StopFilter{}
}

func (f *StopFilter) ByIDs(ids ...int64) *StopFilter {
	f.IDs = ids
	return f
}

func (f *StopFilter) ByCities(cities ...string) *StopFilter {
	f.Cities = cities
	return f
}

func (f *StopFilter) ByAddresses(addresses ...string) *StopFilter {
	f.Addresses = addresses
	return f
}
