package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes/internal/model"
)

type UserStore interface {
	WithTx(*Tx) UserStore
	GetByFilter(ctx context.Context, filter *UserFilter) (*model.User, error)
	GetListByFilter(ctx context.Context, filter *UserFilter) ([]*model.User, error)
	Add(ctx context.Context, users ...*model.User) (int64, error)
	Delete(ctx context.Context, filter *UserFilter) error
	Update(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, hashedPassword []byte, filter *UserFilter) error
}

type UserFilter struct {
	IDs              []int
	Emails           []string
	DoSelectPassword bool
	DoSelectType     bool
}

func NewUserFilter() *UserFilter {
	return &UserFilter{}
}

// ByIDs filters by user.id.
func (f *UserFilter) ByIDs(ids ...int) *UserFilter {
	f.IDs = ids
	return f
}

// ByEmails filters by user.email.
func (f *UserFilter) ByEmails(emails ...string) *UserFilter {
	f.Emails = emails
	return f
}

// SelectPassword selects user.hash_password.
func (f *UserFilter) SelectPassword() *UserFilter {
	f.DoSelectPassword = true
	return f
}

// SelectType selects user.type.

func (f *UserFilter) SelectType() *UserFilter {
	f.DoSelectType = true
	return f
}
