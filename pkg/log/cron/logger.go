package cron

import (
	"context"

	"go-scaffold/pkg/log"

	"golang.org/x/exp/slog"
)

// Logger adapts to cron.Logger
type Logger struct {
	logger  *log.Logger
	logInfo bool
}

func NewLogger(sl *slog.Logger, logInfo bool) *Logger {
	return &Logger{
		logger:  log.NewLogger(sl),
		logInfo: logInfo,
	}
}

func (l *Logger) Info(msg string, kvs ...any) {
	if l.logInfo {
		l.logger.Log(context.Background(), 3, nil, slog.LevelInfo, msg, kvs...)
	}
}

func (l *Logger) Error(err error, msg string, kvs ...any) {
	l.logger.Log(context.Background(), 3, err, slog.LevelError, msg, kvs...)
}
