package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes/internal/model"
)

type BusStore interface {
	WithTx(*Tx) BusStore
	ByFilter(ctx context.Context, filter *BusFilter) (*model.Bus, error)
	ListByFilter(ctx context.Context, filter *BusFilter) ([]*model.Bus, error)
}

type BusFilter struct {
	IDs []int64
}

func NewBusFilter() *BusFilter {
	return &BusFilter{}
}

func (f *BusFilter) ByIDs(ids ...int64) *BusFilter {
	f.IDs = ids
	return f
}
