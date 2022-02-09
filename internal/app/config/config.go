package config

import (
	restconfig "go-scaffold/internal/app/rest/config"
	"go-scaffold/pkg/logger"
	"go-scaffold/pkg/orm"
	"go-scaffold/pkg/redisclient"
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

type Config struct {
	Name             string
	Env              Env
	ShutdownWaitTime int
	Log              *logger.Config
	DB               *orm.Config
	Redis            *redisclient.Config
	REST             restconfig.Config `mapstructure:"REST"`
}
