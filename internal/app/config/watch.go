package config

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
)

var watchKeys = []string{}

func Watch(hLogger log.Logger, cfg config.Config, cm *Config) error {
	var logger = log.NewHelper(hLogger)

	for _, key := range watchKeys {
		logger.Infof("the config is being watching, key: %s", key)

		if err := cfg.Watch(key, func(s string, value config.Value) {
			logger.Infof("config has changed, key: %s", s)

			if err := cfg.Scan(cm); err != nil {
				logger.Errorf("scan config to model failed, err: %v", err)
			}
		}); err != nil {
			return err
		}
	}

	return nil
}
