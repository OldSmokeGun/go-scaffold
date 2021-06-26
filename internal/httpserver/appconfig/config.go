package appconfig

import (
	"gin-scaffold/pkg/logger"
	"gin-scaffold/pkg/orm"
	"gin-scaffold/pkg/redisclient"
	"time"
)

type (
	Config struct {
		AppConf      `mapstructure:",squash"`
		LogConf      *logger.Config      `mapstructure:"Log"`
		TemplateConf *TemplateConf       `mapstructure:"Template"`
		JwtConf      *JwtConf            `mapstructure:"Jwt"`
		DatabaseConf *orm.Config         `mapstructure:"DB"`
		RedisConf    *redisclient.Config `mapstructure:"Redis"`
	}

	AppConf struct {
		Host string
		Port int
		Env  string
	}

	TemplateConf struct {
		Glob string
	}

	JwtConf struct {
		Key    string
		Expire time.Duration
	}
)
