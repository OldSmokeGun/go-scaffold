package config

import (
	"go-scaffold/internal/app/component/casbin"
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
	App       *App              `json:"app"`
	HTTP      *HTTP             `json:"http"`
	GRPC      *GRPC             `json:"grpc"`
	DB        *orm.Config       `json:"db"`
	Redis     *redis.Config     `json:"redis"`
	Trace     *trace.Config     `json:"trace"`
	Discovery *discovery.Config `json:"discovery"`
	Services  *Services         `json:"services"`
	Jwt       *Jwt              `json:"jwt"`
	Casbin    *casbin.Config    `json:"casbin"`
}

type App struct {
	Name    string `json:"name"`
	Env     Env    `json:"env"`
	Timeout int64  `json:"timeout"`
}

type HTTP struct {
	Network      string `json:"network"`
	Addr         string `json:"addr"`
	Timeout      int64  `json:"timeout"`
	ExternalAddr string `json:"externalAddr"`
}

type GRPC struct {
	Network string `json:"network"`
	Addr    string `json:"addr"`
	Timeout int64  `json:"timeout"`
}

type Services struct {
	Self string `json:"self"`
}

type Jwt struct {
	Key string `json:"key"`
}
