package gorm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-scaffold/pkg/log"

	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config is logger config
type Config struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogInfo                   bool
}

// Logger adapts to gorm logger
type Logger struct {
	logger *log.Logger
	config Config
}

func NewLogger(sl *slog.Logger, config Config) *Logger {
	return &Logger{
		logger: log.NewLogger(sl),
		config: config,
	}
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *Logger) Info(ctx context.Context, s string, i ...interface{}) {
	l.logger.Log(ctx, 5, nil, slog.LevelInfo, s, i...)
}

func (l *Logger) Warn(ctx context.Context, s string, i ...interface{}) {
	l.logger.Log(ctx, 5, nil, slog.LevelWarn, s, i...)
}

func (l *Logger) Error(ctx context.Context, s string, i ...interface{}) {
	l.logger.Log(ctx, 5, nil, slog.LevelError, s, i...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	switch {
	case err != nil && (!l.config.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		l.logger.Log(ctx, 5, err, slog.LevelError, "query error",
			slog.String("elapsed", elapsed.String()),
			slog.Int64("rows", rows),
			slog.String("sql", sql),
		)
	case l.config.SlowThreshold != 0 && elapsed > l.config.SlowThreshold:
		msg := fmt.Sprintf("slow threshold >= %v", l.config.SlowThreshold)
		l.logger.Log(ctx, 5, nil, slog.LevelWarn, msg,
			slog.String("elapsed", elapsed.String()),
			slog.Int64("rows", rows),
			slog.String("sql", sql),
		)
	case l.config.LogInfo:
		l.logger.Log(ctx, 5, nil, slog.LevelInfo, "query info",
			slog.String("elapsed", elapsed.String()),
			slog.Int64("rows", rows),
			slog.String("sql", sql),
		)
	}
}
