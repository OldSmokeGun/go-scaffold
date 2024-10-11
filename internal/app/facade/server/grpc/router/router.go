package router

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"

	v1api "go-scaffold/internal/app/facade/server/grpc/api/v1"
)

// Router 注册器
type Router struct {
	greetServer      v1api.GreetServer
	userServer       v1api.UserServer
	roleServer       v1api.RoleServer
	permissionServer v1api.PermissionServer
	productServer    v1api.ProductServer
}

// New 构造注册器
func New(
	greetServer v1api.GreetServer,
	userServer v1api.UserServer,
	roleServer v1api.RoleServer,
	permissionServer v1api.PermissionServer,
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
	v1api.RegisterUserServer(server, r.userServer)
	v1api.RegisterRoleServer(server, r.roleServer)
	v1api.RegisterPermissionServer(server, r.permissionServer)
	v1api.RegisterProductServer(server, r.productServer)
}
