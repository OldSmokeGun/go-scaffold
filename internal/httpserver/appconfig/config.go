package appconfig

import (
	"gin-scaffold/pkg/logger"
	"gin-scaffold/pkg/orm"
	"gin-scaffold/pkg/redisclient"
	"time"
)

type AppEnv string

const (
	Local      AppEnv = "local"
	Test       AppEnv = "test"
	Production AppEnv = "prod"
)

func (a AppEnv) String() string {
	return string(a)
}

type (
	Config struct {
		AppConf      `mapstructure:",squash"`
		LogConf      LogConf             `mapstructure:"Log"`
		LoggerConf   *logger.Config      `mapstructure:"Logger"`
		TemplateConf *TemplateConf       `mapstructure:"Template"`
		DatabaseConf *orm.Config         `mapstructure:"Database"`
		RedisConf    *redisclient.Config `mapstructure:"Redis"`
		JwtConf      *JwtConf            `mapstructure:"Jwt"`
	}

	AppConf struct {
		Host string
		Port int
		Env  AppEnv
	}

	LogConf struct {
		Path string
	}

	TemplateConf struct {
		Glob string
	}

	JwtConf struct {
		Key    string
		Expire time.Duration
	}
)
