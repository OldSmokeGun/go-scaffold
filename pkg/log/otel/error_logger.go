package otel

import (
	"context"
	"log/slog"

	"go-scaffold/pkg/log"
)

// ErrorLogger adapts to otel ErrorLogger
type ErrorLogger struct {
	logger *log.Logger
}

func NewLogger(sl *slog.Logger) *ErrorLogger {
	return &ErrorLogger{
		logger: log.NewLogger(sl),
	}
}

func (l *ErrorLogger) Handle(err error) {
	l.logger.Log(context.Background(), 7, err, slog.LevelError, "otel error")
}
