package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(context.Context, string, ...zapcore.Field)
	Error(context.Context, string, ...zapcore.Field)
	Info(context.Context, string, ...zapcore.Field)
	Named(string) Logger
	Warn(context.Context, string, ...zapcore.Field)
	WithOptions(...zap.Option) Logger
}

type defaultLogger struct {
	zapLogger *zap.Logger
}

func NewLogger(config zap.Config) (Logger, error) {
	config.DisableStacktrace = true
	config.DisableCaller = true
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	z, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &defaultLogger{
		zapLogger: z,
	}, nil
}

func NewNoopLogger() Logger {
	return &defaultLogger{
		zapLogger: zap.NewNop(),
	}
}

func (l *defaultLogger) Debug(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.zapLogger.Debug(
		msg,
		fields...,
	)
}

func (l *defaultLogger) Error(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.zapLogger.Error(
		msg,
		fields...,
	)
}

func (l *defaultLogger) Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.zapLogger.Info(
		msg,
		fields...,
	)
}

func (l *defaultLogger) Named(name string) Logger {
	zl := l.zapLogger.Named(name)
	return &defaultLogger{
		zapLogger: zl,
	}
}

func (l *defaultLogger) Warn(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.zapLogger.Warn(
		msg,
		fields...,
	)
}

func (l *defaultLogger) WithOptions(options ...zap.Option) Logger {
	zl := l.zapLogger.WithOptions(options...)
	return &defaultLogger{
		zapLogger: zl,
	}
}
