package logger

import (
	"gin-scaffold/internal/global"
	"gin-scaffold/internal/utils"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

// Register 注册全局日志对象
func Register(logPath string) (*logrus.Logger, error) {
	var (
		err       error
		logWriter *os.File
		logger    = logrus.New()
	)

	if logPath != "" {
		if !filepath.IsAbs(logPath) {
			logPath = filepath.Join(filepath.Dir(global.GetBinPath()), logPath)
		}

		if ok := utils.PathExist(logPath); !ok {
			logDir := logPath
			if ok, _ := utils.IsDir(logPath); !ok {
				logDir = filepath.Dir(logPath)
			}
			if err := os.MkdirAll(logDir, 0666); err != nil {
				return nil, err
			}
		}

		logWriter, err = os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
	}

	logger.SetOutput(io.MultiWriter(logWriter, os.Stdout))
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	return logger, nil
}
