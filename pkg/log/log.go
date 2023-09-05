package log

import (
	"context"
	"fmt"
	"go/build"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"

	ierrors "go-scaffold/pkg/errors"
	"go-scaffold/pkg/path"
)

const AttrErrorKey = "error"

var (
	// DefaultLevel default log level
	DefaultLevel = Info

	// DefaultFormat default log format
	DefaultFormat = Json

	// DefaultWriter default Writer
	DefaultWriter = os.Stdout
)

// config is logger config
type config struct {
	level  Level
	format Format
	writer io.Writer
	group  string
	attrs  []slog.Attr
}

type Option func(c *config)

// WithLevel set log level
func WithLevel(level Level) Option {
	return func(c *config) {
		c.level = level
	}
}

// WithFormat set log format
func WithFormat(format Format) Option {
	return func(c *config) {
		c.format = format
	}
}

// WithWriter set log writer
func WithWriter(writer io.Writer) Option {
	return func(c *config) {
		c.writer = writer
	}
}

// WithGroup set log group
func WithGroup(group string) Option {
	return func(c *config) {
		c.group = group
	}
}

// WithAttrs set log key-value pair
func WithAttrs(attrs []slog.Attr) Option {
	return func(c *config) {
		c.attrs = attrs
	}
}

// New build *slog.Logger
func New(options ...Option) *slog.Logger {
	logger := &config{
		level:  DefaultLevel,
		format: DefaultFormat,
		writer: DefaultWriter,
	}

	for _, opf := range options {
		opf(logger)
	}

	ops := &slog.HandlerOptions{
		AddSource:   true,
		Level:       logger.level.Convert(),
		ReplaceAttr: ReplaceAttr,
	}

	var handler slog.Handler
	if logger.format == Text {
		handler = slog.NewTextHandler(logger.writer, ops)
	} else {
		handler = slog.NewJSONHandler(logger.writer, ops)
	}

	if logger.group != "" {
		handler = handler.WithGroup(logger.group)
	}

	if len(logger.attrs) > 0 {
		handler = handler.WithAttrs(logger.attrs)
	}

	return slog.New(handler)
}

// Format log format
type Format string

const (
	Text Format = "text"
	Json Format = "json"
)

// Level log level
type Level string

const (
	Debug Level = "debug"
	Info  Level = "info"
	Warn  Level = "warn"
	Error Level = "error"
)

// Convert from Level to slog.Level
func (l Level) Convert() slog.Level {
	switch l {
	case Debug:
		return slog.LevelDebug
	case Info:
		return slog.LevelInfo
	case Warn:
		return slog.LevelWarn
	case Error:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// ReplaceAttr handle log key-value pair
func ReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		return slog.String(a.Key, a.Value.Time().Format(time.RFC3339))
	case slog.LevelKey:
		return slog.String(a.Key, strings.ToLower(a.Value.String()))
	case slog.SourceKey:
		value := a.Value.Any().(*slog.Source)
		return slog.String(a.Key, fmt.Sprintf("%s:%d", getBriefSource(value.File), value.Line))
	case AttrErrorKey:
		v, ok := a.Value.Any().(interface {
			StackTrace() errors.StackTrace
		})
		if ok {
			st := v.StackTrace()
			return slog.Any(a.Key, slog.GroupValue(slog.String("msg", a.Value.String()), slog.Any("stack", ierrors.StackTrace(st))))
		}
		return a
	}
	return a
}

// NewNop returns a no-op logger
func NewNop() *slog.Logger {
	nopLevel := slog.Level(-99)
	ops := &slog.HandlerOptions{
		Level: nopLevel,
	}
	handler := slog.NewTextHandler(io.Discard, ops)
	return slog.New(handler)
}

// Logger is wrapper of slog
type Logger struct {
	*slog.Logger
}

// NewLogger return wrapper of slog
func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{logger}
}

// Log send log records with caller depth
func (l *Logger) Log(ctx context.Context, depth int, err error, level slog.Level, msg string, attrs ...any) {
	if !l.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(depth, pcs[:])
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	if err != nil {
		r.Add(AttrErrorKey, err)
	}
	r.Add(attrs...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.Handler().Handle(ctx, r)
}

func getBriefSource(source string) string {
	gp := filepath.ToSlash(build.Default.GOPATH)
	if strings.HasPrefix(source, gp) {
		return strings.TrimPrefix(source, gp+"/pkg/mod/")
	}
	pp := filepath.ToSlash(path.ProjectPath())
	return strings.TrimPrefix(source, pp+"/")
}
