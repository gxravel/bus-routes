package busroutescontext

import (
	"context"

	"github.com/gxravel/bus-routes/internal/model"
)

type ctxKey string

const (
	UserTypesKey ctxKey = "user_types"
	UserKey      ctxKey = "user"
	TokenKey     ctxKey = "token"
)

// GetUserTypes returns registered user types.
func GetUserTypes(ctx context.Context) model.UserTypes {
	if ctx == nil {
		return nil
	}

	t, _ := ctx.Value(UserTypesKey).(model.UserTypes)

	return t
}
