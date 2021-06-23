package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes/internal/model"
)

type StopStore interface {
	WithTx(*Tx) StopStore
	ByFilter(ctx context.Context, filter *StopFilter) (*model.Stop, error)
	ListByFilter(ctx context.Context, filter *StopFilter) ([]*model.Stop, error)
	Add(ctx context.Context, stops ...*model.Stop) error
	Update(ctx context.Context, stop *model.Stop) error
	Delete(ctx context.Context, filter *StopFilter) error
}

type StopFilter struct {
	IDs         []int64
	Cities      []string
	CitiesIDs   []int
	Addresses   []string
	DoPreferIDs bool
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

func (f *StopFilter) ByCitiesIDs(citiesIDs ...int) *StopFilter {
	f.CitiesIDs = citiesIDs
	return f
}

func (f *StopFilter) PreferIDs() *StopFilter {
	f.DoPreferIDs = true
	return f
}

func (f *StopFilter) ByAddresses(addresses ...string) *StopFilter {
	f.Addresses = addresses
	return f
}
