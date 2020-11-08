package components

import (
	"gin-scaffold/internal/global"
	"gin-scaffold/internal/utils"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

var DefaultLog = filepath.Join(filepath.Dir(global.BinPath), "../logs/errors.log")

func LoadLogrus(f string) error {
	var (
		logPath = DefaultLog
	)

	if f != "" {
		logPath = f
	}

	if !filepath.IsAbs(logPath) {
		logPath = filepath.Join(filepath.Dir(global.BinPath), logPath)
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

	logrus.SetOutput(io.MultiWriter(logWriter, os.Stdout))
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	return nil
}
