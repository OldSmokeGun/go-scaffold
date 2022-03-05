package grpc

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	greetpb "go-scaffold/internal/app/api/v1/greet"
	userpb "go-scaffold/internal/app/api/v1/user"
	"go-scaffold/internal/app/config"
	"time"
)

var ProviderSet = wire.NewSet(NewServer)

// NewServer 创建 gRPC 服务器
func NewServer(
	logger log.Logger,
	cm *config.Config,
	greetService greetpb.GreetServer,
	userService userpb.UserServer,
) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(recovery.WithLogger(logger)),
			logging.Server(logger),
			tracing.Server(),
			validate.Validator(),
		),
		grpc.Logger(logger),
	}

	if cm.App.Grpc.Network != "" {
		opts = append(opts, grpc.Network(cm.App.Grpc.Network))
	}

	if cm.App.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(cm.App.Grpc.Addr))
	}

	if cm.App.Grpc.Timeout != 0 {
		opts = append(opts, grpc.Timeout(time.Duration(cm.App.Grpc.Timeout)*time.Second))
	}

	srv := grpc.NewServer(opts...)

	greetpb.RegisterGreetServer(srv, greetService)
	userpb.RegisterUserServer(srv, userService)

	return srv
}
