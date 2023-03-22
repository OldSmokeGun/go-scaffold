package config

import (
	kconfig "github.com/go-kratos/kratos/v2/config"
	"golang.org/x/exp/slog"
)

var watchKeys = []string{
	"services.self",
}

// Watch 监听配置键的变化
func Watch(logger *slog.Logger, source kconfig.Config, cm *Config) error {
	logger.With(slog.Any("keys", watchKeys)).Debug("the config is being watching")

	for _, key := range watchKeys {
		if err := source.Watch(key, func(s string, value kconfig.Value) {
			logger.With(slog.String("key", key)).Debug("config has changed")

			if err := source.Scan(cm); err != nil {
				logger.Error("scan config to model failed", err)
			}
		}); err != nil {
			return err
		}
	}
	return nil
}
