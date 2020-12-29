package global

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
)

var (
	err error
	// 二进制文件路径
	binPath string
	// 全局 DB
	db           *gorm.DB
	configurator *viper.Viper
	logger       *logrus.Logger
)

func init() {
	binPath, err = os.Executable()
	if err != nil {
		panic(err)
	}
}

// BinPath 获取二进制文件路径
func BinPath() string {
	return binPath
}

// SetDB 设置全局数据库操作对象
func SetDB(d *gorm.DB) {
	db = d
}

// DB 获取全局数据库操作对象
func DB() *gorm.DB {
	return db
}

// SetConfigurator 设置全局配置对象
func SetConfigurator(c *viper.Viper) {
	configurator = c
}

// Configurator 获取全局配置对象
func Configurator() *viper.Viper {
	return configurator
}

// SetLogger 设置全局日志对象
func SetLogger(l *logrus.Logger) {
	logger = l
}

// Logger 获取全局日志对象
func Logger() *logrus.Logger {
	return logger
}
