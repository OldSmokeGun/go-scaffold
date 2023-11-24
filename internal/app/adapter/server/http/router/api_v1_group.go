package router

import (
	"github.com/labstack/echo/v4"

	v1 "go-scaffold/internal/app/adapter/server/http/handler/v1"
	imiddleware "go-scaffold/internal/app/adapter/server/http/middleware"
	"go-scaffold/internal/app/controller"
)

// ApiV1Group v1 API routing group
type ApiV1Group struct {
	accountTokenController      *controller.AccountTokenController
	accountPermissionController *controller.AccountPermissionController

	greetHandler      *v1.GreetHandler
	traceHandler      *v1.TraceHandler
	producerHandler   *v1.ProducerHandler
	accountHandler    *v1.AccountHandler
	userHandler       *v1.UserHandler
	roleHandler       *v1.RoleHandler
	permissionHandler *v1.PermissionHandler
	productHandler    *v1.ProductHandler

	group *echo.Group

	basePath string
}

// NewAPIV1Group return *ApiV1Group
func NewAPIV1Group(
	accountTokenController *controller.AccountTokenController,
	accountPermissionController *controller.AccountPermissionController,
	greetHandler *v1.GreetHandler,
	traceHandler *v1.TraceHandler,
	producerHandler *v1.ProducerHandler,
	accountHandler *v1.AccountHandler,
	userHandler *v1.UserHandler,
	roleHandler *v1.RoleHandler,
	permissionHandler *v1.PermissionHandler,
	productHandler *v1.ProductHandler,
) *ApiV1Group {
	return &ApiV1Group{
		accountTokenController:      accountTokenController,
		accountPermissionController: accountPermissionController,
		greetHandler:                greetHandler,
		traceHandler:                traceHandler,
		productHandler:              productHandler,
		accountHandler:              accountHandler,
		userHandler:                 userHandler,
		roleHandler:                 roleHandler,
		permissionHandler:           permissionHandler,
		producerHandler:             producerHandler,
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

	g.group.POST("/register", g.accountHandler.Register)
	g.group.POST("/login", g.accountHandler.Login)

	g.group.Use(imiddleware.Auth(*imiddleware.NewDefaultAuthConfig().
		WithTokenValidator(g.accountTokenController).
		WithTokenRefresher(g.accountTokenController),
	))
	{
		g.group.DELETE("/logout", g.accountHandler.Logout)
		g.group.PUT("/account/profile", g.accountHandler.UpdateProfile)
		g.group.GET("/account/profile", g.accountHandler.GetProfile)
		g.group.GET("/account/permissions", g.accountHandler.GetPermissions)

		g.group.Use(imiddleware.Permission(*imiddleware.NewDefaultPermissionConfig().
			WithValidator(g.accountPermissionController),
		))

		g.group.GET("/users", g.userHandler.List)
		g.group.GET("/user/:id", g.userHandler.Detail)
		g.group.POST("/user", g.userHandler.Create)
		g.group.PUT("/user", g.userHandler.Update)
		g.group.DELETE("/user/:id", g.userHandler.Delete)
		g.group.GET("/user/roles", g.userHandler.GetRoles)
		g.group.POST("/user/roles", g.userHandler.AssignRoles)

		g.group.GET("/roles", g.roleHandler.List)
		g.group.GET("/role/:id", g.roleHandler.Detail)
		g.group.POST("/role", g.roleHandler.Create)
		g.group.PUT("/role", g.roleHandler.Update)
		g.group.DELETE("/role/:id", g.roleHandler.Delete)
		g.group.GET("/role/permissions", g.roleHandler.GetPermissions)
		g.group.POST("/role/permissions", g.roleHandler.GrantPermissions)

		g.group.GET("/permissions", g.permissionHandler.List)
		g.group.GET("/permission/:id", g.permissionHandler.Detail)
		g.group.POST("/permission", g.permissionHandler.Create)
		g.group.PUT("/permission", g.permissionHandler.Update)
		g.group.DELETE("/permission/:id", g.permissionHandler.Delete)

		g.group.GET("/products", g.productHandler.List)
		g.group.GET("/product/:id", g.productHandler.Detail)
		g.group.POST("/product", g.productHandler.Create)
		g.group.PUT("/product", g.productHandler.Update)
		g.group.DELETE("/product/:id", g.productHandler.Delete)
	}
}
