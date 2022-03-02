package config

import (
	"github.com/go-kratos/kratos/v2/config"
	"go-scaffold/internal/app/component/orm"
	"time"
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
	App *App
}

type App struct {
	Name    string
	Env     Env
	Timeout int
	DB      *struct {
		Driver          string
		Host            string
		Port            string
		Database        string
		Username        string
		Password        string
		Options         []string
		MaxIdleConn     int
		MaxOpenConn     int
		ConnMaxLifeTime int64
		LogLevel        orm.LogLevel
	}
	Redis *struct {
		Host               string
		Port               int
		Username           string
		Password           string
		DB                 int
		MaxRetries         int
		MinRetryBackoff    time.Duration
		MaxRetryBackoff    time.Duration
		DialTimeout        time.Duration
		ReadTimeout        time.Duration
		WriteTimeout       time.Duration
		PoolSize           int
		MinIdleConns       int
		MaxConnAge         time.Duration
		PoolTimeout        time.Duration
		IdleTimeout        time.Duration
		IdleCheckFrequency time.Duration
	}
	Trace *struct {
		Endpoint string
	}
	Http *struct {
		Network      string
		Addr         string
		Timeout      int
		ExternalAddr string
	}
	Grpc *struct {
		Network string
		Addr    string
		Timeout int
	}
	Discovery *struct {
		Endpoints []string
	}
	Jwt *struct {
		Key string
	}
}

func Watch(cfg config.Config, model *Config) error {
	// if err := cfg.Watch("key", func(s string, value config.Value) {
	// 	cfg.Scan(model)
	// }); err != nil {
	//
	// }

	return nil
}
