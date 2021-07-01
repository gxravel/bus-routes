package busroutes

import (
	"context"

	v1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes/internal/dataprovider"
	ierr "github.com/gxravel/bus-routes/internal/errors"
	"github.com/gxravel/bus-routes/internal/logger"
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

func (r *BusRoutes) GetUsers(ctx context.Context, filter *dataprovider.UserFilter) ([]*v1.User, error) {
	dbUsers, err := r.userStore.ListByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return users(dbUsers...), nil
}

func (r *BusRoutes) CheckPasswordHash(ctx context.Context, password string, filter *dataprovider.UserFilter) error {
	dbUser, err := r.userStore.ByFilter(ctx, filter)
	if err != nil {
		return err
	}
	if dbUser == nil {
		return ierr.NewReason(ierr.ErrWrongCredentials)
	}
	return checkPasswordHash(password, dbUser.HashedPassword)
}

func (r *BusRoutes) GetUserType(ctx context.Context, filter *dataprovider.UserFilter) (model.UserType, error) {
	dbUser, err := r.userStore.ByFilter(ctx, filter)
	if err != nil {
		return "", err
	}
	if dbUser == nil {
		return "", ierr.NewReason(ierr.ErrWrongCredentials)
	}
	return dbUser.Type, nil
}

func (r *BusRoutes) AddUsers(ctx context.Context, users ...*v1.User) error {
	err := r.userStore.Add(ctx, dbUsers(ctx, users...)...)
	if err != nil {
		err = ierr.CheckDuplicate(err, "email")
		return err
	}
	return nil
}

func (r *BusRoutes) UpdateUser(ctx context.Context, user *v1.User) error {
	return r.userStore.Update(ctx, dbUsers(ctx, user)[0])
}

func (r *BusRoutes) UpdateUserPassword(ctx context.Context, hashedPassword []byte, filter *dataprovider.UserFilter) error {
	return r.userStore.UpdatePassword(ctx, hashedPassword, filter)
}

func (r *BusRoutes) DeleteUser(ctx context.Context, filter *dataprovider.UserFilter) error {
	return r.userStore.Delete(ctx, filter)
}

func dbUsers(ctx context.Context, users ...*v1.User) []*model.User {
	var dbUsers = make([]*model.User, 0, len(users))
	logger := logger.FromContext(ctx)
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

func users(dbUsers ...*model.User) []*v1.User {
	var users = make([]*v1.User, 0, len(dbUsers))
	for _, user := range dbUsers {
		users = append(users, &v1.User{
			ID:    user.ID,
			Email: user.Email,
			Type:  user.Type,
		})
	}
	return users
}
