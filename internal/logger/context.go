package logger

import "context"

// auxilary logger key type.
type loggerContextKey = struct{}

// CtxWithLogger puts target logger into context.
func CtxWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}
