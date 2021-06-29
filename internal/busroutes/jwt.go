package busroutes

import (
	"context"

	v1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
)

func (r *BusRoutes) NewJWT(ctx context.Context, user *v1.User) (*v1.Token, error) {
	return r.jwtManager.SetNew(ctx, user)
}
