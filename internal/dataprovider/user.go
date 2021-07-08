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

func (f *UserFilter) ByIDs(ids ...int) *UserFilter {
	f.IDs = ids
	return f
}

func (f *UserFilter) ByEmails(emails ...string) *UserFilter {
	f.Emails = emails
	return f
}
func (f *UserFilter) SelectPassword() *UserFilter {
	f.DoSelectPassword = true
	return f
}
func (f *UserFilter) SelectType() *UserFilter {
	f.DoSelectType = true
	return f
}
