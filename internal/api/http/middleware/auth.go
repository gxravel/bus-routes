package middleware

import (
	"context"
	"net/http"
	"strings"

	api "github.com/gxravel/bus-routes/internal/api/http"
	"github.com/gxravel/bus-routes/internal/busroutes"
	"github.com/gxravel/bus-routes/internal/busroutescontext"
	"github.com/gxravel/bus-routes/internal/model"
)

// RegisterUserTypes adds to request's context allowed user types.
func RegisterUserTypes(types ...model.UserType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, busroutescontext.UserTypesKey, model.UserTypes(types))
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// Auth searches user by token and adds his data to context.
func Auth(busroutes *busroutes.BusRoutes) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			allowedUserTypes := busroutescontext.GetUserTypes(ctx)
			token := getAuthToken(r)

			user, err := busroutes.UserByToken(ctx, token, allowedUserTypes...)
			if err != nil {
				api.RespondError(ctx, w, err)
				return
			}

			ctx = context.WithValue(ctx, busroutescontext.UserKey, user)
			ctx = context.WithValue(ctx, busroutescontext.TokenKey, token)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

const (
	// AuthHeader is a header used to find token of user.
	AuthHeader = "Authorization"
)

func getAuthToken(r *http.Request) string {
	tokens, ok := r.Header[AuthHeader]
	if ok {
		if len(tokens) == 1 {
			return strings.TrimPrefix(tokens[0], "Bearer ")
		}
	}

	return ""
}
