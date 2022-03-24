package discovery

import (
	"github.com/go-kratos/kratos/v2/registry"
	"go-scaffold/internal/app/component/discovery/consul"
	"go-scaffold/internal/app/component/discovery/etcd"
	"go-scaffold/internal/app/config"
	"go.uber.org/zap"
)

type Discovery interface {
	registry.Registrar
	registry.Discovery
}

type Config struct {
	Etcd   *etcd.Config
	Consul *consul.Config
}

func NewConfig(dcvConfig *config.Discovery) *Config {
	if dcvConfig == nil {
		return nil
	}

	return &Config{
		Etcd: &etcd.Config{
			Endpoints: dcvConfig.Etcd.Endpoints,
		},
		Consul: &consul.Config{
			Addr:   dcvConfig.Consul.Addr,
			Schema: dcvConfig.Consul.Schema,
		},
	}
}

func New(config *Config, logger *zap.Logger) (Discovery, error) {
	if config == nil {
		return nil, nil
	}

	if config.Etcd != nil {
		return etcd.New(config.Etcd, logger)
	}

	if config.Consul != nil {
		return consul.New(config.Consul)
	}

	return nil, nil
}
