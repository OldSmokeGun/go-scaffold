package router

import (
	"github.com/labstack/echo/v4"

	v1 "go-scaffold/internal/app/adapter/server/http/handler/v1"
	imiddleware "go-scaffold/internal/app/adapter/server/http/middleware"
	"go-scaffold/internal/app/controller"
)

// ApiV1Group v1 API routing group
type ApiV1Group struct {
	systemSessionTokenController      *controller.SystemSessionTokenController
	systemSessionPermissionController *controller.SystemSessionPermissionController

	greetHandler         *v1.GreetHandler
	traceHandler         *v1.TraceHandler
	producerHandler      *v1.ProducerHandler
	systemSessionHandler *v1.SystemSessionHandler
	userHandler          *v1.SystemUserHandler
	roleHandler          *v1.SystemRoleHandler
	permissionHandler    *v1.SystemPermissionHandler
	productHandler       *v1.ProductHandler

	group *echo.Group

	basePath string
}

// NewAPIV1Group return *ApiV1Group
func NewAPIV1Group(
	systemSessionTokenController *controller.SystemSessionTokenController,
	systemSessionPermissionController *controller.SystemSessionPermissionController,
	greetHandler *v1.GreetHandler,
	traceHandler *v1.TraceHandler,
	producerHandler *v1.ProducerHandler,
	systemSessionHandler *v1.SystemSessionHandler,
	userHandler *v1.SystemUserHandler,
	roleHandler *v1.SystemRoleHandler,
	permissionHandler *v1.SystemPermissionHandler,
	productHandler *v1.ProductHandler,
) *ApiV1Group {
	return &ApiV1Group{
		systemSessionTokenController:      systemSessionTokenController,
		systemSessionPermissionController: systemSessionPermissionController,
		greetHandler:                      greetHandler,
		traceHandler:                      traceHandler,
		productHandler:                    productHandler,
		systemSessionHandler:              systemSessionHandler,
		userHandler:                       userHandler,
		roleHandler:                       roleHandler,
		permissionHandler:                 permissionHandler,
		producerHandler:                   producerHandler,
	}
}

func (g *ApiV1Group) setup(prefix string, rg *echo.Group) {
	path := "/v1"
	g.group = rg.Group(path)
	g.basePath = prefix + path
}

func (g *ApiV1Group) useRoutes() {
	g.group.GET("/greet", g.greetHandler.Hello)
	g.group.POST("/trace/example", g.traceHandler.Example)
	g.group.POST("/producer/example", g.producerHandler.Example)

	g.group.POST("/register", g.systemSessionHandler.Register)
	g.group.POST("/login", g.systemSessionHandler.Login)

	g.group.Use(imiddleware.Auth(*imiddleware.NewDefaultAuthConfig().
		WithTokenValidator(g.systemSessionTokenController).
		WithTokenRefresher(g.systemSessionTokenController),
	))
	{
		g.group.DELETE("/logout", g.systemSessionHandler.Logout)
		g.group.PUT("/profile", g.systemSessionHandler.UpdateProfile)
		g.group.PUT("/profile/password", g.systemSessionHandler.UpdatePassword)
		g.group.GET("/profile", g.systemSessionHandler.GetProfile)
		g.group.GET("/permissions", g.systemSessionHandler.GetPermissions)

		g.group.Use(imiddleware.Permission(*imiddleware.NewDefaultPermissionConfig().
			WithValidator(g.systemSessionPermissionController),
		))

		g.group.GET("/system-users", g.userHandler.List)
		g.group.GET("/system-user/:id", g.userHandler.Detail)
		g.group.POST("/system-user", g.userHandler.Create)
		g.group.PUT("/system-user", g.userHandler.Update)
		g.group.DELETE("/system-user/:id", g.userHandler.Delete)
		g.group.GET("/system-user/roles", g.userHandler.GetRoles)
		g.group.POST("/system-user/roles", g.userHandler.AssignRoles)

		g.group.GET("/system-roles", g.roleHandler.List)
		g.group.GET("/system-role/:id", g.roleHandler.Detail)
		g.group.POST("/system-role", g.roleHandler.Create)
		g.group.PUT("/system-role", g.roleHandler.Update)
		g.group.DELETE("/system-role/:id", g.roleHandler.Delete)
		g.group.GET("/system-role/permissions", g.roleHandler.GetPermissions)
		g.group.POST("/system-role/permissions", g.roleHandler.GrantPermissions)

		g.group.GET("/system-permissions", g.permissionHandler.List)
		g.group.GET("/system-permission/:id", g.permissionHandler.Detail)
		g.group.POST("/system-permission", g.permissionHandler.Create)
		g.group.PUT("/system-permission", g.permissionHandler.Update)
		g.group.DELETE("/system-permission/:id", g.permissionHandler.Delete)

		g.group.GET("/products", g.productHandler.List)
		g.group.GET("/product/:id", g.productHandler.Detail)
		g.group.POST("/product", g.productHandler.Create)
		g.group.PUT("/product", g.productHandler.Update)
		g.group.DELETE("/product/:id", g.productHandler.Delete)
	}
}
