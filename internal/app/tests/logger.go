package tests

import (
	"go-scaffold/pkg/log"
	"go.uber.org/zap"
	"io"
	"os"
)

// NewLogger 初始化日志
func NewLogger() *zap.Logger {
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
		log.WithCallerSkip(logCallerSkip),
	)

	return logger
}
