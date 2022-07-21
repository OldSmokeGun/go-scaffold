package component

import (
	"go-scaffold/internal/app/component/casbin"
	"go-scaffold/internal/app/component/client/grpc"
	"go-scaffold/internal/app/component/discovery"
	"go-scaffold/internal/app/component/ent"
	"go-scaffold/internal/app/component/orm"
	"go-scaffold/internal/app/component/redis"
	"go-scaffold/internal/app/component/trace"
	"go-scaffold/internal/app/component/uid"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(orm.New),
	wire.NewSet(ent.New),
	wire.NewSet(redis.New),
	wire.NewSet(trace.New),
	wire.NewSet(discovery.New),
	wire.NewSet(casbin.New),
	wire.NewSet(wire.Bind(new(uid.Generator), new(*uid.Uid)), uid.New),
	grpc.New,
)
