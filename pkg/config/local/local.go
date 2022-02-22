package local

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var ErrConfigPathMissing = errors.New("the config path missing")

type Local struct {
	Path string // 本地配置文件路径

	viper *viper.Viper
}

func New(path string) (*Local, error) {
	if path == "" {
		return nil, ErrConfigPathMissing
	}

	v := viper.New()
	v.SetConfigName(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
	v.AddConfigPath(filepath.Dir(path))

	return &Local{
		Path:  path,
		viper: v,
	}, nil
}

// Load 加载配置到模型
func (l *Local) Load(model interface{}) error {

	if err := l.viper.ReadInConfig(); err != nil {
		return err
	}

	if err := l.viper.Unmarshal(model); err != nil {
		return err
	}

	return nil
}

// MustLoad 加载配置到模型
func (l *Local) MustLoad(model interface{}) {
	if err := l.Load(model); err != nil {
		panic(err)
	}
}

// Watch 监控配置变更
func (l *Local) Watch(f func(*viper.Viper, fsnotify.Event)) {
	l.viper.WatchConfig()
	l.viper.OnConfigChange(func(e fsnotify.Event) {
		f(l.viper, e)
	})
}
