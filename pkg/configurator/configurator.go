package configurator

import (
	"errors"
	"gin-scaffold/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrFileNotExist = errors.New("the specified file does not exist")
)

// LoadConfig 加载配置
func LoadConfig(path string, v interface{}) error {
	var cfg = viper.New()

	if path != "" {
		if !filepath.IsAbs(path) {
			path = filepath.Join(filepath.Dir(global.GetBinPath()), path)
		}

		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return ErrFileNotExist
			}
		}

		cfg.SetConfigName(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))

		if filepath.IsAbs(path) {
			cfg.AddConfigPath(filepath.Dir(path))
		} else {
			cfg.AddConfigPath(filepath.Dir(filepath.Join(filepath.Dir(global.GetBinPath()), path)))
		}

		cfg.WatchConfig()

		cfg.OnConfigChange(func(in fsnotify.Event) {
			if err := cfg.MergeInConfig(); err != nil {
				panic(err)
			}
			if err := cfg.Unmarshal(v); err != nil {
				panic(err)
			}
		})

		if err := cfg.ReadInConfig(); err != nil {
			return err
		}

		if err := cfg.Unmarshal(v); err != nil {
			return err
		}
	}

	return nil
}

// MustLoadConfig 加载配置
func MustLoadConfig(path string, v interface{}) {
	if err := LoadConfig(path, v); err != nil {
		panic(err)
	}
}
