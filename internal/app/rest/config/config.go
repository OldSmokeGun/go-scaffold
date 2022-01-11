package config

import (
	"gin-scaffold/pkg/logger"
	"gin-scaffold/pkg/orm"
	"gin-scaffold/pkg/redisclient"
	"io"
	"time"
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
		App   `mapstructure:",squash"`
		Log   *Log                `mapstructure:"Log"`
		DB    *orm.Config         `mapstructure:"DB"`
		Redis *redisclient.Config `mapstructure:"Redis"`
		Jwt   *Jwt                `mapstructure:"Jwt"`
	}

	App struct {
		Host string
		Port int
		Env  Env
	}

	Log struct {
		Path   string
		Level  logger.Level
		Format logger.Format
		Caller bool
		Mode   logger.Mode
		Output io.Writer
	}

	Jwt struct {
		Key    string
		Expire time.Duration
	}
)
