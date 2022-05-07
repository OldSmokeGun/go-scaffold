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
	"go-scaffold/internal/app/transport/grpc/handler"
	"time"
)

var ProviderSet = wire.NewSet(
	handler.ProviderSet,
	NewServer,
)

// NewServer 创建 gRPC 服务器
func NewServer(
	logger log.Logger,
	grpcConf *config.Grpc,
	greetServer greetpb.GreetServer,
	userServer userpb.UserServer,
) *grpc.Server {
	if grpcConf == nil {
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

	if grpcConf.Network != "" {
		opts = append(opts, grpc.Network(grpcConf.Network))
	}

	if grpcConf.Addr != "" {
		opts = append(opts, grpc.Address(grpcConf.Addr))
	}

	if grpcConf.Timeout != 0 {
		opts = append(opts, grpc.Timeout(time.Duration(grpcConf.Timeout)*time.Second))
	}

	srv := grpc.NewServer(opts...)

	greetpb.RegisterGreetServer(srv, greetServer)
	userpb.RegisterUserServer(srv, userServer)

	return srv
}
