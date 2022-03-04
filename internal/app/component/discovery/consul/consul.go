package consul

import (
	consul "github.com/go-kratos/consul/registry"
	"github.com/hashicorp/consul/api"
)

type Config struct {
	Addr   string
	Schema string
}

// New 创建 consul 服务发现
func New(config *Config) (*consul.Registry, error) {
	if config == nil {
		return nil, nil
	}

	defaultConfig := api.DefaultConfig()
	if config.Addr != "" {
		defaultConfig.Address = config.Addr
	}
	if config.Schema != "" {
		defaultConfig.Scheme = config.Schema
	}

	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return nil, err
	}

	return consul.New(client), nil
}
