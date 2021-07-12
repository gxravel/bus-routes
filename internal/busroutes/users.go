package busroutes

import (
	"context"

	httpv1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes/internal/dataprovider"
	ierr "github.com/gxravel/bus-routes/internal/errors"
	log "github.com/gxravel/bus-routes/internal/logger"
	"github.com/gxravel/bus-routes/internal/model"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = bcrypt.DefaultCost
)

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
}

func checkPasswordHash(password string, hashedPassword []byte) error {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return ierr.NewReason(ierr.ErrWrongCredentials)
	}

	return nil
}

func (r *BusRoutes) GetUsers(ctx context.Context, filter *dataprovider.UserFilter) ([]*httpv1.User, error) {
	dbUsers, err := r.userStore.GetListByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	return toV1Users(dbUsers...), nil
}

func (r *BusRoutes) CheckPasswordHash(ctx context.Context, password string, filter *dataprovider.UserFilter) error {
	dbUser, err := r.userStore.GetByFilter(ctx, filter)
	if err != nil {
		return err
	}
	if dbUser == nil {
		return ierr.NewReason(ierr.ErrWrongCredentials)
	}

	return checkPasswordHash(password, dbUser.HashedPassword)
}

func (r *BusRoutes) GetUserType(ctx context.Context, filter *dataprovider.UserFilter) (model.UserType, error) {
	dbUser, err := r.userStore.GetByFilter(ctx, filter)
	if err != nil {
		return "", err
	}
	if dbUser == nil {
		return "", ierr.NewReason(ierr.ErrWrongCredentials)
	}

	return dbUser.Type, nil
}

func (r *BusRoutes) AddUsers(ctx context.Context, users ...*httpv1.User) (int64, error) {
	id, err := r.userStore.Add(ctx, toDBUsers(ctx, users...)...)
	if err != nil {
		err = ierr.CheckDuplicate(err, "email")
		return 0, err
	}

	return id, nil
}

func (r *BusRoutes) UpdateUser(ctx context.Context, user *httpv1.User) error {
	return r.userStore.Update(ctx, toDBUsers(ctx, user)[0])
}

func (r *BusRoutes) UpdateUserPassword(ctx context.Context, hashedPassword []byte, filter *dataprovider.UserFilter) error {
	return r.userStore.UpdatePassword(ctx, hashedPassword, filter)
}

func (r *BusRoutes) DeleteUser(ctx context.Context, filter *dataprovider.UserFilter) error {
	return r.userStore.Delete(ctx, filter)
}

func toDBUsers(ctx context.Context, users ...*httpv1.User) []*model.User {
	var dbUsers = make([]*model.User, 0, len(users))

	logger := log.FromContext(ctx)

	for _, user := range users {
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			logger.Debug(err.Error())
		}

		dbUsers = append(dbUsers, &model.User{
			ID:             user.ID,
			Email:          user.Email,
			Type:           user.Type,
			HashedPassword: hashedPassword,
		})
	}

	return dbUsers
}

func toV1Users(dbUsers ...*model.User) []*httpv1.User {
	var users = make([]*httpv1.User, 0, len(dbUsers))
	for _, user := range dbUsers {
		users = append(users, &httpv1.User{
			ID:    user.ID,
			Email: user.Email,
			Type:  user.Type,
		})
	}

	return users
}
