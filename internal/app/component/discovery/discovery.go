package discovery

import (
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	etcdctl "go.etcd.io/etcd/client/v3"
)

type Config struct {
	Endpoints []string
}

// New 创建 etcd 服务发现
func New(config *Config) (*etcd.Registry, error) {
	if config == nil {
		return nil, nil
	}

	client, err := etcdctl.New(etcdctl.Config{
		Endpoints: config.Endpoints,
	})
	if err != nil {
		return nil, err
	}

	return etcd.New(client), nil
}
