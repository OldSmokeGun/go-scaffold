package http

//go:generate swag fmt -g http.go
//go:generate swag init -g http.go -o api/docs --parseInternal

import (
	"net/http"
	"time"

	v1 "go-scaffold/internal/app/adapter/server/http/handler/v1"
	"go-scaffold/internal/app/adapter/server/http/router"
	"go-scaffold/internal/config"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
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
//	@name						Token

var ProviderSet = wire.NewSet(
	// handler
	v1.NewGreetHandler,
	v1.NewTraceHandler,
	v1.NewUserHandler,
	v1.NewProducerHandler,
	// router
	router.New,
	router.NewAPIGroup,
	router.NewAPIV1Group,
	// HTTP server
	New,
)

// New build HTTP server
func New(
	httpConf config.HTTP,
	handler http.Handler,
) *khttp.Server {
	if handler == nil {
		return nil
	}

	var opts []khttp.ServerOption

	if httpConf.Network != "" {
		opts = append(opts, khttp.Network(httpConf.Network))
	}

	if httpConf.Addr != "" {
		opts = append(opts, khttp.Address(httpConf.Addr))
	}

	if httpConf.Timeout != 0 {
		opts = append(opts, khttp.Timeout(httpConf.Timeout*time.Second))
	}

	srv := khttp.NewServer(opts...)
	srv.HandlePrefix("/", handler)

	return srv
}
