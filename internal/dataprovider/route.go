package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes/internal/model"
)

type RouteStore interface {
	WithTx(*Tx) RouteStore
	GetByFilter(ctx context.Context, filter *RouteFilter) (*model.Route, error)
	GetListByFilter(ctx context.Context, filter *RouteFilter) ([]*model.Route, error)
	Add(ctx context.Context, routes ...*model.Route) error
	Update(ctx context.Context, route *model.Route) error
	Delete(ctx context.Context, filter *RouteFilter) error
}

type RouteFilter struct {
	BusIDs       []int64
	StopIDs      []int64
	Steps        []int8
	DetailedView bool
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

func (f *RouteFilter) ViewDetailed() *RouteFilter {
	f.DetailedView = true
	return f
}
