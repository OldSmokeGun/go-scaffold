package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

var logger = logrus.New()

// Setup 返回 *logrus.Logger
func Setup(conf Config) (*logrus.Logger, error) {
	var err error

	path := conf.Path

	if path == "" {
		logger.SetOutput(io.MultiWriter(conf.Output, os.Stdout))
	} else {
		// 如果路径不存在，则创建
		_, err = os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				dir := path
				if filepath.Ext(path) != "" {
					dir = filepath.Dir(path)
				}
				if err := os.MkdirAll(dir, 0666); err != nil {
					return nil, err
				}
			}
		}

		logWriter, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}

		logger.SetOutput(io.MultiWriter(logWriter, os.Stdout))
	}

	logger.SetLevel(conf.Level.Convert())
	logger.SetReportCaller(conf.ReportCaller)

	switch conf.Format {
	case Text:
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000",
		})
	case Json:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000",
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000",
		})
	}

	return logger, nil
}

// MustSetup 返回 *logrus.Logger
func MustSetup(conf Config) *logrus.Logger {
	l, err := Setup(conf)
	if err != nil {
		panic(err)
	}

	return l
}
