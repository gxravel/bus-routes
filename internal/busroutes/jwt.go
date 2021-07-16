package busroutes

import (
	"context"

	httpv1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	ierr "github.com/gxravel/bus-routes/internal/errors"
	"github.com/gxravel/bus-routes/internal/jwt"
	log "github.com/gxravel/bus-routes/internal/logger"
	"github.com/gxravel/bus-routes/internal/model"
)

func (r *Busroutes) NewJWT(ctx context.Context, user *httpv1.User) (*httpv1.Token, error) {
	jwtUser := &jwt.User{
		ID:    user.ID,
		Email: user.Email,
		Type:  user.Type,
	}
	details, err := r.tokenManager.SetNew(ctx, jwtUser)
	if err != nil {
		return nil, err
	}
	return &httpv1.Token{
		Token:  details.String,
		Expiry: details.Expiry,
	}, nil
}

// GetUserByToken returns user withdrawn from the JWT token claims, unless it is of not allowed type.
func (r *Busroutes) GetUserByToken(ctx context.Context, token string, allowedTypes ...model.UserType) (*httpv1.User, error) {
	logger := log.FromContext(ctx).WithStr("token", token)

	if token == "" {
		err := ierr.NewReason(ierr.ErrInvalidToken)
		return nil, err
	}

	jwtUser, err := r.tokenManager.Verify(ctx, token)
	if err != nil {
		return nil, err
	}

	user := &httpv1.User{
		ID:    jwtUser.ID,
		Email: jwtUser.Email,
		Type:  jwtUser.Type,
	}

	types := model.UserTypes(allowedTypes)
	if !types.Exists(user.Type) {
		err := ierr.NewReason(ierr.ErrPermissionDenied)

		logger.
			WithField("user", user).
			Warn(err.Error())

		return nil, err
	}

	return user, nil
}
