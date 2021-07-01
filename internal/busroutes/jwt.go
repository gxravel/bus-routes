package busroutes

import (
	"context"

	v1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	ierr "github.com/gxravel/bus-routes/internal/errors"
	"github.com/gxravel/bus-routes/internal/logger"
	"github.com/gxravel/bus-routes/internal/model"
)

func (r *BusRoutes) NewJWT(ctx context.Context, user *v1.User) (*v1.Token, error) {
	return r.tokenManager.SetNew(ctx, user)
}

func (r *BusRoutes) GetUserByToken(ctx context.Context, token string, allowedTypes ...model.UserType) (*v1.User, error) {
	logger := logger.FromContext(ctx).WithStr("token", token)

	if token == "" {
		err := ierr.NewReason(ierr.ErrInvalidToken)
		return nil, err
	}

	user, err := r.tokenManager.Verify(ctx, token)
	if err != nil {
		return nil, err
	}

	types := model.UserTypes(allowedTypes)
	if !types.Exists(user.Type) {
		err := ierr.NewReason(ierr.ErrPermissionDenied)
		logger.WithField("user", user).Warn(err.Error())
		return nil, err
	}
	return user, nil
}
