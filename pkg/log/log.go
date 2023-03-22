package log

import (
	"context"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	ierrors "go-scaffold/pkg/errors"
	"go-scaffold/pkg/path"

	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

var (
	// DefaultLevel default log level
	DefaultLevel = Info

	// DefaultFormat default log format
	DefaultFormat = Json

	// DefaultWriter default Writer
	DefaultWriter = os.Stdout
)

// option is logger option
type option struct {
	level  Level
	format Format
	writer io.Writer
	group  string
	attrs  []slog.Attr
}

// OptionFunc optional function
type OptionFunc func(logger *option)

// WithLevel set log level
func WithLevel(level Level) OptionFunc {
	return func(logger *option) {
		logger.level = level
	}
}

// WithFormat set log format
func WithFormat(format Format) OptionFunc {
	return func(logger *option) {
		logger.format = format
	}
}

// WithWriter set log writer
func WithWriter(writer io.Writer) OptionFunc {
	return func(logger *option) {
		logger.writer = writer
	}
}

// WithGroup set log group
func WithGroup(group string) OptionFunc {
	return func(logger *option) {
		logger.group = group
	}
}

// WithAttrs set log key-value pair
func WithAttrs(attrs []slog.Attr) OptionFunc {
	return func(logger *option) {
		logger.attrs = attrs
	}
}

// New build *slog.Logger
func New(options ...OptionFunc) *slog.Logger {
	logger := &option{
		level:  DefaultLevel,
		format: DefaultFormat,
		writer: DefaultWriter,
	}

	for _, opf := range options {
		opf(logger)
	}

	ops := slog.HandlerOptions{
		AddSource:   true,
		Level:       logger.level.Convert(),
		ReplaceAttr: ReplaceAttr,
	}

	var handler slog.Handler
	if logger.format == Text {
		handler = ops.NewTextHandler(logger.writer)
	} else {
		handler = ops.NewJSONHandler(logger.writer)
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
		return slog.String(a.Key, getBriefSource(a.Value.String()))
	case slog.ErrorKey:
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
	ops := slog.HandlerOptions{
		Level: nopLevel,
	}
	handler := ops.NewTextHandler(io.Discard)
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
		r.Add(slog.ErrorKey, err)
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
