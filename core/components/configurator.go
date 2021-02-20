package components

import (
	"gin-scaffold/core/global"
	"gin-scaffold/core/utils"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var DefaultConfigPath = filepath.Join(filepath.Dir(filepath.Dir(global.BinPath())), "config/config.yaml")

// RegisterConfigurator 注册全局配置对象
func RegisterConfigurator(configPath string) error {
	var configurator = viper.New()

	if configPath != "" {
		if ok := utils.PathExist(configPath); ok {
			configurator.SetConfigName(strings.TrimSuffix(filepath.Base(configPath), filepath.Ext(configPath)))

			if filepath.IsAbs(configPath) {
				configurator.AddConfigPath(filepath.Dir(configPath))
			} else {
				configurator.AddConfigPath(filepath.Dir(filepath.Join(filepath.Dir(global.BinPath()), configPath)))
			}

			configurator.WatchConfig()

			if err := configurator.ReadInConfig(); err != nil {
				return err
			}
		}
	}

	// 设置全局配置对象
	global.SetConfigurator(configurator)

	return nil
}
