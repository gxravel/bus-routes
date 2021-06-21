package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gxravel/bus-routes/internal/logger"
)

func Logger(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = logger.CtxWithLogger(ctx, log)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.FromContext(r.Context()).Errorf("panic: %v", err)
				debug.PrintStack()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
