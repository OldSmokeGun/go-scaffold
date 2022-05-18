package config

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"go-scaffold/internal/app/component/casbin"
	"go-scaffold/internal/app/component/orm"
	"go-scaffold/internal/app/model"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var ProviderSet = wire.NewSet(
	wire.FieldsOf(
		new(*Config),
		"App",
		"HTTP",
		"GRPC",
		"DB",
		"Redis",
		"Trace",
		"Discovery",
		"Services",
		"JWT",
		"Casbin",
	),
	wire.FieldsOf(new(*DB), "Config"),
	wire.FieldsOf(new(*casbin.Config), "Model", "Adapter"),
)

var watchKeys = []string{
	"services.self",
}

// Loaded 配置加载后调用的钩子函数
func Loaded(hLogger log.Logger, cfg config.Config, conf *Config) error {
	if conf.Trace != nil {
		conf.Trace.ServiceName = conf.App.Name
		conf.Trace.Env = conf.App.Env.String()
		conf.Trace.Timeout = conf.App.Timeout
	}

	if conf.Casbin != nil {
		if conf.Casbin.Adapter != nil {
			if conf.Casbin.Adapter.Gorm != nil {
				conf.Casbin.Adapter.Gorm.SetMigration(func(db *gorm.DB) error {
					return (&model.CasbinRule{}).Migrate(db)
				})
			}
		}
	}

	// 配置多数据库
	if conf.DB != nil {
		if len(conf.DB.Resolvers) > 0 {
			var (
				sources  = make([]gorm.Dialector, 0, len(conf.DB.Resolvers))
				replicas = make([]gorm.Dialector, 0, len(conf.DB.Resolvers))
			)

			for _, resolver := range conf.DB.Resolvers {
				dial, err := orm.BuildDialector(conf.DB.Driver, resolver.DSN)
				if err != nil {
					return err
				}
				switch resolver.Type {
				case orm.Source:
					sources = append(sources, dial)
				case orm.Replica:
					replicas = append(replicas, dial)
				default:
					return fmt.Errorf("unsupported resolver type %s", resolver.Type)
				}
			}

			conf.DB.Plugins = func(db *gorm.DB) ([]gorm.Plugin, error) {
				return []gorm.Plugin{
					dbresolver.Register(dbresolver.Config{
						Sources:  sources,
						Replicas: replicas,
						Policy:   dbresolver.RandomPolicy{},
					}),
				}, nil
			}
		}
	}

	if err := Watch(hLogger, cfg, conf); err != nil {
		return err
	}

	return nil
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
