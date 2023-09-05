package tests

import (
	"io"
	"log/slog"
	"os"

	"go-scaffold/pkg/log"
)

// NewLogger init slog logger
func NewLogger() *slog.Logger {
	var writer io.Writer
	if logLevel == "silent" {
		writer = io.Discard
	} else {
		writer = os.Stdout
	}

	logger := log.New(
		log.WithLevel(log.Level(logLevel)),
		log.WithFormat(log.Format(logFormat)),
		log.WithWriter(writer),
		log.WithAttrs([]slog.Attr{
			slog.String("app", appName),
		}),
	)

	return logger
}
