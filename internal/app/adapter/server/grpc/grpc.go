package grpc

import (
	"time"

	v1api "go-scaffold/internal/app/adapter/server/grpc/api/v1"
	v1handler "go-scaffold/internal/app/adapter/server/grpc/handler/v1"
	"go-scaffold/internal/app/adapter/server/grpc/router"
	"go-scaffold/internal/config"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	// handler
	wire.NewSet(wire.Bind(new(v1api.GreetServer), new(*v1handler.GreetHandler)), v1handler.NewGreetHandler),
	wire.NewSet(wire.Bind(new(v1api.UserServer), new(*v1handler.UserHandler)), v1handler.NewUserHandler),
	// register
	router.New,
	// gRPC server
	New,
)

// New build gRPC server
func New(
	grpcConf config.GRPC,
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

	if grpcConf.Network != "" {
		opts = append(opts, grpc.Network(grpcConf.Network))
	}

	if grpcConf.Addr != "" {
		opts = append(opts, grpc.Address(grpcConf.Addr))
	}

	if grpcConf.Timeout != 0 {
		opts = append(opts, grpc.Timeout(grpcConf.Timeout*time.Second))
	}

	srv := grpc.NewServer(opts...)

	router.Register(srv)

	return srv
}
