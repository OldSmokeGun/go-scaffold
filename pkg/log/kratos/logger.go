package kratos

import (
	"context"

	"go-scaffold/pkg/log"

	klog "github.com/go-kratos/kratos/v2/log"
	"golang.org/x/exp/slog"
)

var _ klog.Logger = (*Logger)(nil)

type Logger struct {
	logger *log.Logger
}

func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{log.NewLogger(logger)}
}

func (l *Logger) Log(level klog.Level, kvs ...interface{}) error {
	argLen := len(kvs)

	if argLen <= 1 {
		return nil
	}
	if kvs[0] != klog.DefaultMessageKey {
		return nil
	}
	msg, ok := kvs[1].(string)
	if !ok {
		return nil
	}

	var attrs []any
	for i := 2; i < argLen; i++ {
		attrs = append(attrs, kvs[i])
	}

	ctx := context.Background()
	switch level {
	case klog.LevelDebug:
		l.logger.Log(ctx, 4, nil, slog.LevelDebug, msg, attrs...)
	case klog.LevelInfo:
		l.logger.Log(ctx, 4, nil, slog.LevelInfo, msg, attrs...)
	case klog.LevelWarn:
		l.logger.Log(ctx, 4, nil, slog.LevelWarn, msg, attrs...)
	case klog.LevelError:
		l.logger.Log(ctx, 4, nil, slog.LevelError, msg, attrs...)
	}
	return nil
}

func (l *Logger) Close() error { return nil }
