package configurator

import (
	"errors"
	"gin-scaffold/internal/global"
	"gin-scaffold/internal/utils"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var (
	ErrFileNotExist = errors.New("the specified file does not exist")
)

// Register 注册全局配置对象
func Register(configPath string) (*viper.Viper, error) {
	var configurator = viper.New()

	if configPath != "" {
		if !filepath.IsAbs(configPath) {
			configPath = filepath.Join(filepath.Dir(global.GetBinPath()), configPath)
		}

		if ok := utils.PathExist(configPath); !ok {
			return nil, ErrFileNotExist
		}

		configurator.SetConfigName(strings.TrimSuffix(filepath.Base(configPath), filepath.Ext(configPath)))

		if filepath.IsAbs(configPath) {
			configurator.AddConfigPath(filepath.Dir(configPath))
		} else {
			configurator.AddConfigPath(filepath.Dir(filepath.Join(filepath.Dir(global.GetBinPath()), configPath)))
		}

		configurator.WatchConfig()

		if err := configurator.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	return configurator, nil
}
