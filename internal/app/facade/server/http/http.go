package http

//go:generate swag fmt -g http.go
//go:generate swag init -g http.go -o api/docs --parseInternal

import (
	"net/http"
	"time"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	v1 "go-scaffold/internal/app/facade/server/http/handler/v1"
	"go-scaffold/internal/app/facade/server/http/router"
	"go-scaffold/internal/config"
)

//	@title						API 接口文档
//	@description				API 接口文档
//	@version					0.0.0
//	@host						localhost
//	@BasePath					/api
//	@schemes					http https
//	@accept						json
//	@accept						x-www-form-urlencoded
//	@produce					json
//	@securityDefinitions.apikey	Authorization
//	@in							header
//	@name						Authorization

var ProviderSet = wire.NewSet(
	// handler
	v1.NewGreetHandler,
	v1.NewTraceHandler,
	v1.NewProducerHandler,
	v1.NewAccountHandler,
	v1.NewUserHandler,
	v1.NewRoleHandler,
	v1.NewPermissionHandler,
	v1.NewProductHandler,
	// router
	router.New,
	router.NewAPIGroup,
	router.NewAPIV1Group,
	// HTTP server
	New,
)

// New build HTTP server
func New(
	hsConf config.HTTPServer,
	handler http.Handler,
) *khttp.Server {
	if handler == nil {
		return nil
	}

	var opts []khttp.ServerOption

	if hsConf.Network != "" {
		opts = append(opts, khttp.Network(hsConf.Network))
	}

	if hsConf.Addr != "" {
		opts = append(opts, khttp.Address(hsConf.Addr))
	}

	if hsConf.Timeout != 0 {
		opts = append(opts, khttp.Timeout(hsConf.Timeout*time.Second))
	}

	srv := khttp.NewServer(opts...)
	srv.HandlePrefix("/", handler)

	return srv
}
