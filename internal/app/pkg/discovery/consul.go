package discovery

import (
	"go-scaffold/internal/config"

	consul "github.com/go-kratos/consul/registry"
	"github.com/hashicorp/consul/api"
)

// NewConsul build consul discovery
func NewConsul(conf config.Consul) (*consul.Registry, error) {
	defaultConfig := api.DefaultConfig()
	if conf.Addr != "" {
		defaultConfig.Address = conf.Addr
	}
	if conf.Schema != "" {
		defaultConfig.Scheme = conf.Schema
	}

	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return nil, err
	}

	return consul.New(client), nil
}
