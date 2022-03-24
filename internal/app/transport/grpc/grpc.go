package grpc

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	greetpb "go-scaffold/internal/app/api/scaffold/v1/greet"
	userpb "go-scaffold/internal/app/api/scaffold/v1/user"
	"go-scaffold/internal/app/config"
	"time"
)

var ProviderSet = wire.NewSet(NewServer)

// NewServer 创建 gRPC 服务器
func NewServer(
	logger log.Logger,
	conf *config.Config,
	greetService greetpb.GreetServer,
	userService userpb.UserServer,
) *grpc.Server {
	if conf.App.Grpc == nil {
		return nil
	}

	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(recovery.WithLogger(logger)),
			logging.Server(logger),
			tracing.Server(),
			validate.Validator(),
			metadata.Server(),
		),
		grpc.Logger(logger),
	}

	if conf.App.Grpc.Network != "" {
		opts = append(opts, grpc.Network(conf.App.Grpc.Network))
	}

	if conf.App.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(conf.App.Grpc.Addr))
	}

	if conf.App.Grpc.Timeout != 0 {
		opts = append(opts, grpc.Timeout(time.Duration(conf.App.Grpc.Timeout)*time.Second))
	}

	srv := grpc.NewServer(opts...)

	greetpb.RegisterGreetServer(srv, greetService)
	userpb.RegisterUserServer(srv, userService)

	return srv
}
