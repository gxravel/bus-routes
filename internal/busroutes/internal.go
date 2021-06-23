package busroutes

import "context"

func (r *BusRoutes) IsHealthy(ctx context.Context) error {
	return r.db.StatusCheck(ctx)
}
