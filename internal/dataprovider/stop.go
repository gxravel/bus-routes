package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes/internal/model"
)

type StopStore interface {
	WithTx(*Tx) StopStore
	GetByFilter(ctx context.Context, filter *StopFilter) (*model.Stop, error)
	GetListByFilter(ctx context.Context, filter *StopFilter) ([]*model.Stop, error)
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

// ByIDs filters by stop.id.
func (f *StopFilter) ByIDs(ids ...int64) *StopFilter {
	f.IDs = ids
	return f
}

// ByCities filters by city.name.
func (f *StopFilter) ByCities(cities ...string) *StopFilter {
	f.Cities = cities
	return f
}

// ByCitiesIDs filters by stop.city_id.
func (f *StopFilter) ByCitiesIDs(citiesIDs ...int) *StopFilter {
	f.CitiesIDs = citiesIDs
	return f
}

// ByAddresses filters by stop.address.
func (f *StopFilter) ByAddresses(addresses ...string) *StopFilter {
	f.Addresses = addresses
	return f
}

// PreferIDs select ids instead of joined values.
func (f *StopFilter) PreferIDs() *StopFilter {
	f.DoPreferIDs = true
	return f
}
