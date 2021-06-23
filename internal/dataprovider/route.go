package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes/internal/model"
)

type RouteStore interface {
	WithTx(*Tx) RouteStore
	ByFilter(ctx context.Context, filter *RouteFilter) (*model.Route, error)
	ListByFilter(ctx context.Context, filter *RouteFilter) ([]*model.Route, error)
	New(ctx context.Context, routes ...*model.Route) error
	Update(ctx context.Context, route *model.Route) error
	Delete(ctx context.Context, filter *RouteFilter) error
}

type RouteFilter struct {
	BusIDs  []int64
	StopIDs []int64
	Steps   []int8
}

func NewRouteFilter() *RouteFilter {
	return &RouteFilter{}
}

func (f *RouteFilter) ByBusIDs(ids ...int64) *RouteFilter {
	f.BusIDs = ids
	return f
}

func (f *RouteFilter) ByStopIDs(ids ...int64) *RouteFilter {
	f.StopIDs = ids
	return f
}

func (f *RouteFilter) BySteps(steps ...int8) *RouteFilter {
	f.Steps = steps
	return f
}
