package consul

import (
	consul "github.com/go-kratos/consul/registry"
	"github.com/hashicorp/consul/api"
)

type Config struct {
	Address string
	Scheme  string
}

// New 创建 consul 服务发现
func New(config *Config) (*consul.Registry, error) {
	if config == nil {
		return nil, nil
	}

	defaultConfig := api.DefaultConfig()
	if config.Address != "" {
		defaultConfig.Address = config.Address
	}
	if config.Scheme != "" {
		defaultConfig.Scheme = config.Scheme
	}

	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return nil, err
	}

	return consul.New(client), nil
}
