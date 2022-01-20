package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

// New 返回 *zap.Logger
func New(conf Config) (logger *zap.Logger, err error) {
	var (
		encoderConfig zapcore.EncoderConfig
		encoder       zapcore.Encoder
		writeSyncer   zapcore.WriteSyncer
		stackLevel    zapcore.Level
	)

	switch conf.Mode {
	case Development:
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		stackLevel = zapcore.WarnLevel
	case Production:
		encoderConfig = zap.NewProductionEncoderConfig()
		stackLevel = zapcore.ErrorLevel
	default:
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	switch conf.Format {
	case Text:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case Json:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	if conf.Output == nil {
		writeSyncer = zapcore.AddSync(os.Stderr)
	} else {
		writeSyncer = zapcore.AddSync(io.MultiWriter(conf.Output, os.Stderr))
	}

	core := zapcore.NewCore(
		encoder,
		writeSyncer,
		conf.Level.Convert(),
	)

	if conf.Caller {
		logger = zap.New(
			core,
			zap.AddStacktrace(stackLevel),
			zap.AddCaller(),
		)
	} else {
		logger = zap.New(
			core,
			zap.AddStacktrace(stackLevel),
		)
	}

	return logger, nil
}

// MustNew 返回 *zap.Logger
func MustNew(conf Config) *zap.Logger {
	l, err := New(conf)
	if err != nil {
		panic(err)
	}

	return l
}
