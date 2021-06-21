package logger

import "context"

// auxilary logger key type.
type loggerContextKey = struct{}

// FromContext extract logger from context or return
// new logger instance with default parameters.
func FromContext(ctx context.Context) Logger {
	if ctx != nil {
		if v, ok := ctx.Value(loggerContextKey{}).(Logger); ok {
			return v
		}
	}

	return Default()
}

// CtxWithLogger puts target logger into context.
func CtxWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}
