package ctxlogger

import (
	"context"

	"go.uber.org/zap"
)

// Define a key type to avoid key collisions in context.
type contextKey string

const loggerKey contextKey = "logger"

func SetLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		logger, _ = zap.NewProduction()
		return logger
	}
	return logger
}
