package router

import (
	v1api "go-scaffold/internal/app/adapter/server/grpc/api/v1"

	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// Router 注册器
type Router struct {
	greetServer v1api.GreetServer
	userServer  v1api.UserServer
}

// New 构造注册器
func New(
	greetServer v1api.GreetServer,
	userServer v1api.UserServer,
) *Router {
	return &Router{
		greetServer: greetServer,
		userServer:  userServer,
	}
}

// Register 注册服务
func (r *Router) Register(server *grpc.Server) {
	v1api.RegisterGreetServer(server, r.greetServer)
	v1api.RegisterUserServer(server, r.userServer)
}
