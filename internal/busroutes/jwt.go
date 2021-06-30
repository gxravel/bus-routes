package busroutes

import (
	"context"
	"errors"

	v1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes/internal/logger"
	"github.com/gxravel/bus-routes/internal/model"
)

var (
	ErrTokenNotFound    = errors.New("token not found")
	ErrPermissionDenied = errors.New("user does not have permission")
	ErrTokenExpired     = errors.New("token expired")
)

func (r *BusRoutes) NewJWT(ctx context.Context, user *v1.User) (*v1.Token, error) {
	return r.jwtManager.SetNew(ctx, user)
}

func (r *BusRoutes) UserByToken(ctx context.Context, token string, allowedTypes ...model.UserType) (*v1.User, error) {
	logger := logger.FromContext(ctx).WithStr("token", token)

	if token == "" {
		err := ErrTokenNotFound
		logger.Debug(err.Error())
		return nil, err
	}

	user, err := r.jwtManager.Verify(ctx, token)
	if err != nil {
		logger.WithErr(err).Warn("verifying token")
		return nil, err
	}

	types := model.UserTypes(allowedTypes)
	if !types.Exists(user.Type) {
		err := ErrPermissionDenied
		logger.WithField("user", user).Warn(err.Error())
		return nil, err
	}
	return user, nil
}
