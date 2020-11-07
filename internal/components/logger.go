package components

import (
	"fmt"
	"gin-scaffold/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
)

const DefaultLog = "logs/errors.log"

func InitLogger() error {
	var (
		logPath = DefaultLog
	)

	if viper.IsSet("errors_log") && viper.GetString("errors_log") != "" {
		logPath = viper.GetString("errors_log")
	}

	if ok, _ := utils.PathExist(logPath); !ok {
		logDir := logPath
		if ok, _ := utils.IsDir(logPath); !ok {
			logDir = filepath.Dir(logPath)
		}
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("创建日志目录 %s 出错，错误信息：%w", logDir, err)
		}
	}

	logWriter, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return fmt.Errorf("打开文件 %s 出错，错误信息：%w", logPath, err)
	}

	logrus.SetOutput(io.MultiWriter(logWriter, os.Stdout))
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	return nil
}
