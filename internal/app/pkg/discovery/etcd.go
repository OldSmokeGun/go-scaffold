package discovery

import (
	"context"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	etcdctl "go.etcd.io/etcd/client/v3"

	"go-scaffold/internal/config"
)

// NewEtcd build etcd discovery
func NewEtcd(ctx context.Context, conf config.Etcd) (*etcd.Registry, error) {
	client, err := etcdctl.New(etcdctl.Config{
		Endpoints: conf.Endpoints,
		Context:   ctx,
	})
	if err != nil {
		return nil, err
	}

	return etcd.New(client), nil
}
