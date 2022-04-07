package config

import (
	"go-scaffold/internal/app/component/discovery"
	"go-scaffold/internal/app/component/orm"
	"go-scaffold/internal/app/component/redis"
	"go-scaffold/internal/app/component/trace"
)

// SupportedEnvs 支持的环境
var SupportedEnvs = []string{Local.String(), Test.String(), Prod.String()}

// Env 当前运行环境
type Env string

func (e Env) String() string {
	return string(e)
}

const (
	Local Env = "local"
	Test  Env = "test"
	Prod  Env = "prod"
)

type Config struct {
	App      *App
	Services *Services
	Jwt      *Jwt
}

type (
	App struct {
		Name      string
		Env       Env
		Timeout   int64
		Http      *Http
		Grpc      *Grpc
		DB        *orm.Config
		Redis     *redis.Config
		Trace     *trace.Config
		Discovery *discovery.Config
	}

	Http struct {
		Network      string
		Addr         string
		Timeout      int64
		ExternalAddr string
	}

	Grpc struct {
		Network string
		Addr    string
		Timeout int64
	}
)

type Services struct {
	Self string
}

type Jwt struct {
	Key string
}
