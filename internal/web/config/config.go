package config

import (
	"gin-scaffold/pkg/logger"
	"gin-scaffold/pkg/orm"
	"gin-scaffold/pkg/redisclient"
	"time"
)

type Env string

const (
	Local      Env = "local"
	Test       Env = "test"
	Production Env = "prod"
)

func (a Env) String() string {
	return string(a)
}

type (
	Config struct {
		App    `mapstructure:",squash"`
		Log    Log                 `mapstructure:"Log"`
		Logger *logger.Config      `mapstructure:"Logger"`
		DB     *orm.Config         `mapstructure:"DB"`
		Redis  *redisclient.Config `mapstructure:"Redis"`
		Jwt    *Jwt                `mapstructure:"Jwt"`
	}

	App struct {
		Host string
		Port int
		Env  Env
	}

	Log struct {
		Path string
	}

	Jwt struct {
		Key    string
		Expire time.Duration
	}
)
