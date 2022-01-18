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

// Config 配置
type Config struct {

	// Path 配置文件路径
	Path string

	// Model 配置文件加载到的模型
	Model interface{}

	// OnConfigChange 配置文件发生改变时的回调函数
	OnConfigChange OnConfigChange

	// Viper viper 实例
	Viper *viper.Viper
}

// New 构造函数
func New(path string, model interface{}, ops ...OptionFunc) *Config {
	config := &Config{
		Path:           path,
		Model:          model,
		OnConfigChange: nil,
	}

	for _, op := range ops {
		op(config)
	}

	return config
}

// OnConfigChange 配置文件发生改变时的回调函数
type OnConfigChange func(*Config, fsnotify.Event)

// OptionFunc 选项函数
type OptionFunc func(config *Config)

// WithOnConfigChange 设置 config 的 OnConfigChange 属性
func WithOnConfigChange(f OnConfigChange) OptionFunc {
	return func(config *Config) {
		config.OnConfigChange = f
	}
}

// Load 加载配置
func (c *Config) Load() error {
	var v = viper.New()

	if c.Path == "" {
		return ErrFileNotSpecified
	}

	v.SetConfigName(strings.TrimSuffix(filepath.Base(c.Path), filepath.Ext(c.Path)))

	v.AddConfigPath(filepath.Dir(c.Path))

	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		c.OnConfigChange(c, e)
	})

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(c.Model); err != nil {
		return err
	}

	c.Viper = v

	return nil
}
