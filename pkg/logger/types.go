package logger

import (
	"go.uber.org/zap/zapcore"
	"io"
)

type Format string

const (
	Text Format = "text"
	Json Format = "json"
)

type Level string

const (
	Debug  Level = "debug"
	Info   Level = "info"
	Warn   Level = "warn"
	Error  Level = "error"
	DPanic Level = "dpanic"
	Panic  Level = "panic"
	Fatal  Level = "fatal"
)

func (l Level) Convert() zapcore.Level {
	switch l {
	case Debug:
		return zapcore.DebugLevel
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Error:
		return zapcore.ErrorLevel
	case DPanic:
		return zapcore.DPanicLevel
	case Panic:
		return zapcore.PanicLevel
	case Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

type Mode string

const (
	Development Mode = "development"
	Production  Mode = "production"
)

type Config struct {
	Path   string
	Level  Level
	Format Format
	Caller bool
	Mode   Mode
	Output io.Writer
}
