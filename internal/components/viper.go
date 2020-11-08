package components

import (
	"gin-scaffold/internal/global"
	"gin-scaffold/internal/utils"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var DefaultConfig = filepath.Join(filepath.Dir(global.BinPath), "config.yaml")

func LoadViper(f string) error {
	var (
		config = DefaultConfig
	)

	if f != "" {
		config = f
	}

	if ok := utils.PathExist(config); !ok {
		return nil
	}

	viper.SetConfigName(strings.TrimSuffix(filepath.Base(config), filepath.Ext(config)))

	if filepath.IsAbs(config) {
		viper.AddConfigPath(filepath.Dir(config))
	} else {
		viper.AddConfigPath(filepath.Dir(filepath.Join(filepath.Dir(global.BinPath), config)))
	}

	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
