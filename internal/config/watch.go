package config

import (
	"log/slog"

	kconfig "github.com/go-kratos/kratos/v2/config"
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
				logger.Error("scan config to model failed", slog.Any("error", err))
			}
		}); err != nil {
			return err
		}
	}
	return nil
}
