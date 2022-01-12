package config

import (
	restconfig "gin-scaffold/internal/app/rest/config"
	"gin-scaffold/pkg/logger"
	"gin-scaffold/pkg/orm"
	"gin-scaffold/pkg/redisclient"
)

type Env string

const (
	Local Env = "local"
	Test  Env = "test"
	Prod  Env = "prod"
)

func (a Env) String() string {
	return string(a)
}

type (
	Config struct {
		App  `mapstructure:",squash"`
		REST restconfig.Config `mapstructure:"REST"`
	}

	App struct {
		Env              Env
		ShutdownWaitTime int
		Log              *logger.Config
		DB               *orm.Config
		Redis            *redisclient.Config
	}
)
