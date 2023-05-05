package ent

import (
	"context"

	"go-scaffold/pkg/log"

	"golang.org/x/exp/slog"
)

const defaultMessageKey = "msg"

// Logger adapts to cron.Logger
type Logger struct {
	logger *log.Logger
}

func NewLogger(sl *slog.Logger) *Logger {
	return &Logger{
		logger: log.NewLogger(sl),
	}
}

func (l *Logger) Log(i ...any) {
	attrs := make([]any, 0, len(i))
	if len(i) == 1 {
		attrs = append(attrs, defaultMessageKey)
	}
	attrs = append(attrs, i...)

	if l.logger != nil {
		l.logger.Log(context.Background(), 14, nil, slog.LevelDebug, "ent debug", attrs...)
	}
}
