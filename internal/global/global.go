package global

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	err          error
	binPath      string         // 二进制文件路径
	rootCommand  *cobra.Command // 根命令
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
func GetBinPath() string {
	return binPath
}

// SetRootCommand 设置根命令
func SetRootCommand(cmd *cobra.Command) {
	rootCommand = cmd
}

// GetRootCommand 获取根命令
func GetRootCommand() *cobra.Command {
	return rootCommand
}

// SetConfigurator 设置全局配置对象
func SetConfigurator(c *viper.Viper) {
	configurator = c
}

// GetConfigurator 获取全局配置对象
func GetConfigurator() *viper.Viper {
	return configurator
}

// SetLogger 设置全局日志对象
func SetLogger(l *logrus.Logger) {
	logger = l
}

// GetLogger 获取全局日志对象
func GetLogger() *logrus.Logger {
	return logger
}
