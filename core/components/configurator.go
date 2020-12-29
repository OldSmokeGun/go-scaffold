package components

import (
	"gin-scaffold/app/util"
	"gin-scaffold/core/global"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var DefaultConfigPath = filepath.Join(filepath.Dir(filepath.Dir(global.BinPath())), "config/config.yaml")

// RegisterConfigurator 注册全局配置对象
func RegisterConfigurator(f string) error {
	var (
		config       = DefaultConfigPath
		configurator = viper.New()
	)

	if f != "" {
		config = f
	}

	if ok := util.PathExist(config); ok {
		configurator.SetConfigName(strings.TrimSuffix(filepath.Base(config), filepath.Ext(config)))

		if filepath.IsAbs(config) {
			configurator.AddConfigPath(filepath.Dir(config))
		} else {
			configurator.AddConfigPath(filepath.Dir(filepath.Join(filepath.Dir(global.BinPath()), config)))
		}

		configurator.WatchConfig()

		if err := configurator.ReadInConfig(); err != nil {
			return err
		}
	}

	// 设置全局配置对象
	global.SetConfigurator(configurator)

	return nil
}
