package router

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"

	v1api "go-scaffold/internal/app/adapter/server/grpc/api/v1"
)

// Router 注册器
type Router struct {
	greetServer      v1api.GreetServer
	userServer       v1api.SystemUserServer
	roleServer       v1api.SystemRoleServer
	permissionServer v1api.SystemPermissionServer
	productServer    v1api.ProductServer
}

// New 构造注册器
func New(
	greetServer v1api.GreetServer,
	userServer v1api.SystemUserServer,
	roleServer v1api.SystemRoleServer,
	permissionServer v1api.SystemPermissionServer,
	productServer v1api.ProductServer,
) *Router {
	return &Router{
		greetServer:      greetServer,
		userServer:       userServer,
		roleServer:       roleServer,
		permissionServer: permissionServer,
		productServer:    productServer,
	}
}

// Register 注册服务
func (r *Router) Register(server *grpc.Server) {
	v1api.RegisterGreetServer(server, r.greetServer)
	v1api.RegisterSystemUserServer(server, r.userServer)
	v1api.RegisterSystemRoleServer(server, r.roleServer)
	v1api.RegisterSystemPermissionServer(server, r.permissionServer)
	v1api.RegisterProductServer(server, r.productServer)
}
