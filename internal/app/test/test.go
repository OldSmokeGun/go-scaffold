package test

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

var (
	// logger 测试用 logger
	// 注意：日志内容将会被丢弃，可用于替换实际的 logger
	logger *zap.Logger
)

func Init() {
	var err error

	logger, err = zap.NewDevelopment(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(io.Discard),
			zapcore.ErrorLevel,
		)
	}))
	if err != nil {
		panic(err)
	}
}

// Logger 测试用 logger
func Logger() *zap.Logger {
	return logger
}
