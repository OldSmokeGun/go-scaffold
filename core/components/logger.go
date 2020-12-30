package components

import (
	"gin-scaffold/core/global"
	"gin-scaffold/core/utils"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

var DefaultLogPath = filepath.Join(filepath.Dir(filepath.Dir(global.BinPath())), "logs/errors.log")

// RegisterLogger 注册全局日志对象
func RegisterLogger(f string) error {
	var (
		logPath = DefaultLogPath
		logger  = logrus.New()
	)

	if f != "" {
		logPath = f
	}

	if !filepath.IsAbs(logPath) {
		logPath = filepath.Join(filepath.Dir(global.BinPath()), logPath)
	}

	if ok := utils.PathExist(logPath); !ok {
		logDir := logPath
		if ok, _ := utils.IsDir(logPath); !ok {
			logDir = filepath.Dir(logPath)
		}
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
	}

	logWriter, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return err
	}

	logger.SetOutput(io.MultiWriter(logWriter, os.Stdout))
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	// 设置全局日志对象
	global.SetLogger(logger)

	return nil
}
