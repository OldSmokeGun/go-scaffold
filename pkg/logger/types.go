package logger

import (
	"github.com/sirupsen/logrus"
	"io"
)

type Format string

const (
	Text Format = "text"
	Json Format = "json"
)

type Level string

const (
	Trace Level = "trace"
	Debug Level = "debug"
	Info  Level = "info"
	Warn  Level = "warn"
	Error Level = "error"
	Fatal Level = "fatal"
	Panic Level = "panic"
)

func (l Level) Convert() logrus.Level {
	switch l {
	case Trace:
		return logrus.TraceLevel
	case Debug:
		return logrus.DebugLevel
	case Info:
		return logrus.InfoLevel
	case Warn:
		return logrus.WarnLevel
	case Error:
		return logrus.ErrorLevel
	case Fatal:
		return logrus.FatalLevel
	case Panic:
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}

type Config struct {
	Path         string
	Level        Level
	Format       Format
	ReportCaller bool
	Output       io.Writer
}
