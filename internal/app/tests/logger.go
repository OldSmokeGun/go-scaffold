package tests

import (
	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go-scaffold/pkg/log"
	"go.uber.org/zap"
	"io"
	"os"
)

// NewZapLogger 初始化 zap logger
func NewZapLogger() *zap.Logger {
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

// NewLogger 初始化日志
func NewLogger(zLogger *zap.Logger) klog.Logger {
	logger := klog.With(
		kzap.NewLogger(zLogger),
		"service.id", hostname,
		"service.name", appName,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)

	return logger
}
