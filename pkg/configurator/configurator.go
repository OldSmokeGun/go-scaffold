package configurator

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrFileNotExist = errors.New("the specified file does not exist")
)

// Load 加载配置
func Load(path string, v interface{}) error {
	var cfg = viper.New()

	if path != "" {
		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return ErrFileNotExist
			}
		}

		cfg.SetConfigName(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))

		cfg.AddConfigPath(filepath.Dir(path))

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

// MustLoad 加载配置
func MustLoad(path string, v interface{}) {
	if err := Load(path, v); err != nil {
		panic(err)
	}
}
