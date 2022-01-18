package config

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var (
	ErrFileNotSpecified = errors.New("the file not specified")
)

// OnConfigChange 配置文件发生改变时的回调函数
type OnConfigChange func(*viper.Viper, interface{}, fsnotify.Event)

// Load 加载配置
func Load(path string, model interface{}, onConfigChange ...OnConfigChange) error {
	var v = viper.New()

	if path == "" {
		return ErrFileNotSpecified
	}

	v.SetConfigName(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))

	v.AddConfigPath(filepath.Dir(path))

	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		for _, f := range onConfigChange {
			f(v, model, e)
		}
	})

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(model); err != nil {
		return err
	}

	return nil
}

// MustLoad 加载配置
func MustLoad(path string, model interface{}, onConfigChange ...OnConfigChange) {
	if err := Load(path, model, onConfigChange...); err != nil {
		panic(err)
	}
}
