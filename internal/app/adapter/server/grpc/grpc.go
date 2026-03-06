package grpc

import (
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"

	v1api "go-scaffold/internal/app/adapter/server/grpc/api/v1"
	v1handler "go-scaffold/internal/app/adapter/server/grpc/handler/v1"
	"go-scaffold/internal/app/adapter/server/grpc/router"
	"go-scaffold/internal/config"
)

var ProviderSet = wire.NewSet(
	// handler
	wire.NewSet(wire.Bind(new(v1api.GreetServer), new(*v1handler.GreetHandler)), v1handler.NewGreetHandler),
	wire.NewSet(wire.Bind(new(v1api.UserServer), new(*v1handler.UserHandler)), v1handler.NewUserHandler),
	wire.NewSet(wire.Bind(new(v1api.RoleServer), new(*v1handler.RoleHandler)), v1handler.NewRoleHandler),
	wire.NewSet(wire.Bind(new(v1api.PermissionServer), new(*v1handler.PermissionHandler)), v1handler.NewPermissionHandler),
	wire.NewSet(wire.Bind(new(v1api.ProductServer), new(*v1handler.ProductHandler)), v1handler.NewProductHandler),
	// register
	router.New,
	// gRPC server
	New,
)

// New build gRPC server
func New(
	gsConf config.GRPCServer,
	router *router.Router,
) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(log.GetLogger()),
			tracing.Server(),
			metadata.Server(),
		),
	}

	if gsConf.Network != "" {
		opts = append(opts, grpc.Network(gsConf.Network))
	}

	if gsConf.Addr != "" {
		opts = append(opts, grpc.Address(gsConf.Addr))
	}

	if gsConf.Timeout != 0 {
		opts = append(opts, grpc.Timeout(gsConf.Timeout*time.Second))
	}

	srv := grpc.NewServer(opts...)

	router.Register(srv)

	return srv
}
