package logger

import (
	"context"
	"errors"
	"syscall"

	"go.uber.org/zap"
)

type contextKey string

const (
	loggerTraceIDKey   contextKey = "x-trace-id"
	loggerRequestIDKey contextKey = "x-request_id"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Sync() error
}

type L struct {
	z *zap.Logger
}

func NewLogger(env string) Logger {
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if env == "dev" {
		loggerCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	logger, err := loggerCfg.Build()
	if err != nil {
		return nil
	}

	return &L{z: logger}
}

func (l *L) Sync() error {
	err := l.z.Sync()
	if err == nil {
		return nil
	}

	if errors.Is(err, syscall.EINVAL) {
		return nil
	}

	return err
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, loggerRequestIDKey, requestID)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, loggerTraceIDKey, traceID)
}

func (l *L) Info(ctx context.Context, msg string, fields ...zap.Field) {
	val := ctx.Value(loggerRequestIDKey)

	if id, ok := val.(string); ok && id != "" {
		fields = append(fields, zap.String(string(loggerRequestIDKey), id))
	}

	l.z.Info(msg, fields...)
}

func (l *L) Error(ctx context.Context, msg string, fields ...zap.Field) {
	val := ctx.Value(loggerRequestIDKey)

	if id, ok := val.(string); ok && id != "" {
		fields = append(fields, zap.String(string(loggerRequestIDKey), id))
	}

	l.z.Error(msg, fields...)
}

func (l *L) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	val := ctx.Value(loggerRequestIDKey)

	if id, ok := val.(string); ok && id != "" {
		fields = append(fields, zap.String(string(loggerRequestIDKey), id))
	}

	l.z.Debug(msg, fields...)
}
