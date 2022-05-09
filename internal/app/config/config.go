package config

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.FieldsOf(new(*Config), "App", "Jwt", "Services"),
	wire.FieldsOf(new(*App), "DB"),
	wire.FieldsOf(new(*App), "Redis"),
	wire.FieldsOf(new(*App), "Trace"),
	wire.FieldsOf(new(*App), "Http"),
	wire.FieldsOf(new(*App), "Grpc"),
	wire.FieldsOf(new(*App), "Discovery"),
)

var watchKeys = []string{
	"services.self",
	"jwt.key",
}

// Watch 监听配置键的变化
func Watch(hLogger log.Logger, cfg config.Config, conf *Config) error {
	var logger = log.NewHelper(hLogger)

	for _, key := range watchKeys {
		logger.Infof("the config is being watching, key: %s", key)

		if err := cfg.Watch(key, func(s string, value config.Value) {
			logger.Infof("config has changed, key: %s", s)

			if err := cfg.Scan(conf); err != nil {
				logger.Errorf("scan config to model failed, err: %v", err)
			}
		}); err != nil {
			return err
		}
	}

	return nil
}

// AfterLoad 配置加载后调用的钩子函数
func AfterLoad(hLogger log.Logger, cfg config.Config, conf *Config) error {
	if conf.App.Trace != nil {
		conf.App.Trace.ServiceName = conf.App.Name
		conf.App.Trace.Env = conf.App.Env.String()
		conf.App.Trace.Timeout = conf.App.Timeout
	}

	if err := Watch(hLogger, cfg, conf); err != nil {
		return err
	}
	return nil
}
