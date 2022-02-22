package config

import (
	restconfig "go-scaffold/internal/app/rest/config"
	"go-scaffold/pkg/orm"
	"go-scaffold/pkg/redisclient"
)

type Config struct {
	Priority         bool
	Name             string
	Env              string
	ShutdownWaitTime int
	DB               *orm.Config
	Redis            *redisclient.Config
	Trace            *struct {
		Endpoint         string
		ShutdownWaitTime int
	}
	REST restconfig.Config `mapstructure:"REST"`
}
