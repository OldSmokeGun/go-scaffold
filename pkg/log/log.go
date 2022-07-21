package log

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// DefaultLevel 默认等级
	DefaultLevel = Info

	// DefaultFormat 默认输出格式
	DefaultFormat = Json

	// DefaultWriter 默认输出 Writer
	DefaultWriter = os.Stdout

	// DefaultCallerSkip 默认跳过 caller 的层级
	DefaultCallerSkip = 0
)

// Logger 日志
type Logger struct {
	Level      Level
	Format     Format
	Writer     io.Writer
	CallerSkip int
}

type Option func(logger *Logger)

// WithLevel 设置日志等级
func WithLevel(level Level) Option {
	return func(logger *Logger) {
		logger.Level = level
	}
}

// WithFormat 设置日志格式
func WithFormat(format Format) Option {
	return func(logger *Logger) {
		logger.Format = format
	}
}

// WithWriter 设置日志输出 writer
func WithWriter(writer io.Writer) Option {
	return func(logger *Logger) {
		logger.Writer = writer
	}
}

// WithCallerSkip 设置日志 caller 跳过的层级
func WithCallerSkip(skip int) Option {
	return func(logger *Logger) {
		logger.CallerSkip = skip
	}
}

// New 返回 *zap.Logger
func New(options ...Option) *zap.Logger {
	logger := &Logger{
		Level:      DefaultLevel,
		Format:     DefaultFormat,
		Writer:     DefaultWriter,
		CallerSkip: DefaultCallerSkip,
	}

	for _, option := range options {
		option(logger)
	}

	var (
		encoderConfig zapcore.EncoderConfig
		encoder       zapcore.Encoder
		writeSyncer   zapcore.WriteSyncer
	)

	encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	switch logger.Format {
	case Text:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case Json:
		fallthrough
	default:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	writeSyncer = zapcore.AddSync(logger.Writer)

	core := zapcore.NewCore(
		encoder,
		writeSyncer,
		logger.Level.Convert(),
	)

	zapLogger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(logger.CallerSkip),
	)

	return zapLogger
}

// Format 日志格式
type Format string

const (
	Text Format = "text"
	Json Format = "json"
)

// Level 日志等级
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

// Convert 转换 Level 为 zapcore.Level
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
