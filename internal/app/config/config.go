package config

import (
	"github.com/google/wire"
	"time"
)

var ProviderSet = wire.NewSet(
	wire.FieldsOf(new(*Config), "App", "Services"),
	wire.FieldsOf(new(*App), "DB"),
	wire.FieldsOf(new(*App), "Redis"),
	wire.FieldsOf(new(*App), "Trace"),
	wire.FieldsOf(new(*App), "Http"),
	wire.FieldsOf(new(*App), "Grpc"),
	wire.FieldsOf(new(*App), "Discovery"),
	wire.FieldsOf(new(*App), "Jwt"),
)

type Config struct {
	App      *App
	Services *Services
}

type (
	App struct {
		Name      string
		Env       Env
		Timeout   int64
		DB        *DB
		Redis     *Redis
		Trace     *Trace
		Http      *Http
		Grpc      *Grpc
		Discovery *Discovery
		Jwt       *Jwt
	}

	DB struct {
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
		LogLevel        string
	}

	Redis struct {
		Host               string
		Port               string
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

	Trace struct {
		Endpoint string
	}

	Http struct {
		Network      string
		Addr         string
		Timeout      int
		ExternalAddr string
	}

	Grpc struct {
		Network string
		Addr    string
		Timeout int
	}

	Discovery struct {
		Etcd *struct {
			Endpoints []string
		}
		Consul *struct {
			Addr   string
			Schema string
		}
	}

	Jwt struct {
		Key string
	}
)

type Services struct {
	Self string
}
