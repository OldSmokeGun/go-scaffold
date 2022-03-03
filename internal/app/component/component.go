package component

import (
	"github.com/google/wire"
	"go-scaffold/internal/app/component/data"
	"go-scaffold/internal/app/component/discovery/consul"
	"go-scaffold/internal/app/component/discovery/etcd"
	"go-scaffold/internal/app/component/orm"
	"go-scaffold/internal/app/component/redis"
	"go-scaffold/internal/app/component/trace"
)

var ProviderSet = wire.NewSet(
	orm.New,
	redis.New,
	trace.New,
	etcd.New,
	consul.New,
	data.New,
)
