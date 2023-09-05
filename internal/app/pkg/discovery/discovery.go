package discovery

import (
	"context"

	"github.com/go-kratos/kratos/v2/registry"

	"go-scaffold/internal/config"
)

// Discovery the interface that discovery must implement
type Discovery interface {
	registry.Registrar
	registry.Discovery
}

// New build service discovery
func New(ctx context.Context, conf config.Discovery) (Discovery, error) {
	if conf.Etcd != nil {
		return NewEtcd(ctx, *conf.Etcd)
	}

	if conf.Consul != nil {
		return NewConsul(*conf.Consul)
	}

	return nil, nil
}
