package mysql

import (
	"context"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type UserStore struct {
	db        sqlx.ExtContext
	txer      dataprovider.Txer
	tableName string
}

func NewUserStore(db sqlx.ExtContext, txer dataprovider.Txer) *UserStore {
	return &UserStore{
		db:        db,
		txer:      txer,
		tableName: "user",
	}
}

func (s *UserStore) WithTx(tx *dataprovider.Tx) dataprovider.UserStore {
	return &UserStore{
		db:        tx,
		tableName: s.tableName,
	}
}

func userCond(f *dataprovider.UserFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	var cond sq.Sqlizer = eq

	if len(f.IDs) > 0 {
		eq["user.id"] = f.IDs
	}

	if len(f.Emails) > 0 {
		eq["email"] = f.Emails
	}

	return cond
}

func (s *UserStore) columns(filter *dataprovider.UserFilter) []string {
	if filter == nil {
		return []string{
			"email",
			"hashed_password",
			"type",
		}
	}
	if filter.SelectPassword {
		return []string{
			"hashed_password",
		}
	}
	return []string{
		"id",
		"email",
		"type",
	}
}

func (s *UserStore) ByFilter(ctx context.Context, filter *dataprovider.UserFilter) (*model.User, error) {
	users, err := s.ListByFilter(ctx, filter)

	switch {
	case err != nil:
		return nil, err
	case len(users) == 0:
		return nil, nil
	case len(users) == 1:
		return users[0], nil
	default:
		return nil, errors.New("fetched more than 1 user")
	}
}

func (s *UserStore) ListByFilter(ctx context.Context, filter *dataprovider.UserFilter) ([]*model.User, error) {
	qb := sq.
		Select(s.columns(filter)...).
		From(s.tableName).
		Where(userCond(filter))

	result, err := selectContext(ctx, qb, s.tableName, s.db, TypeUser)
	if err != nil {
		return nil, err
	}
	return result.([]*model.User), nil
}

func (s *UserStore) Add(ctx context.Context, users ...*model.User) error {
	qb := sq.Insert(s.tableName).Columns(s.columns(nil)...)
	for _, user := range users {
		qb = qb.Values(user.Email, user.HashedPassword, user.Type)
	}
	return execContext(ctx, qb, s.tableName, s.txer)
}

func (s *UserStore) Delete(ctx context.Context, filter *dataprovider.UserFilter) error {
	qb := sq.Delete(s.tableName).Where(userCond(filter))
	return execContext(ctx, qb, s.tableName, s.txer)
}

func (s *UserStore) Update(ctx context.Context, user *model.User) error {
	qb := sq.Update(s.tableName).Set("email", user.Email).Set("type", user.Type).Where(sq.Eq{"id": user.ID})
	return execContext(ctx, qb, s.tableName, s.txer)
}

func (s *UserStore) UpdatePassword(ctx context.Context, hashedPassword []byte, filter *dataprovider.UserFilter) error {
	qb := sq.Update(s.tableName).Set("hashed_password", hashedPassword).Where(userCond(filter))
	return execContext(ctx, qb, s.tableName, s.txer)
}
