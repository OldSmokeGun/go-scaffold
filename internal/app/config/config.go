package config

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
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

func Watch(hLogger log.Logger, cfg config.Config, cm *Config) error {
	var logger = log.NewHelper(hLogger)

	for _, key := range watchKeys {
		logger.Infof("the config is being watching, key: %s", key)

		if err := cfg.Watch(key, func(s string, value config.Value) {
			logger.Infof("config has changed, key: %s", s)

			if err := cfg.Scan(cm); err != nil {
				logger.Errorf("scan config to model failed, err: %v", err)
			}
		}); err != nil {
			return err
		}
	}

	return nil
}

type Config struct {
	App      *App
	Services *Services
}

type App struct {
	Name    string
	Env     Env
	Timeout int64
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
		ConnMaxIdleTime int64
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
		Etcd *struct {
			Endpoints []string
		}
		Consul *struct {
			Addr   string
			Schema string
		}
	}
	Jwt *struct {
		Key string
	}
}

type Services struct {
	Self string
}
