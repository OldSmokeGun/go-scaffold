package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/transport/http/handler"
	"go-scaffold/internal/app/transport/http/router"
	"time"
)

var ProviderSet = wire.NewSet(
	handler.ProviderSet,
	router.ProviderSet,
	NewServer,
)

// NewServer 创建 HTTP 服务器
func NewServer(
	logger log.Logger,
	conf *config.Config,
	router *gin.Engine,
) *khttp.Server {
	if router == nil {
		return nil
	}

	var opts = []khttp.ServerOption{
		khttp.Logger(logger),
	}

	if conf.App.Http.Network != "" {
		opts = append(opts, khttp.Network(conf.App.Http.Network))
	}

	if conf.App.Http.Addr != "" {
		opts = append(opts, khttp.Address(conf.App.Http.Addr))
	}

	if conf.App.Http.Timeout != 0 {
		opts = append(opts, khttp.Timeout(time.Duration(conf.App.Http.Timeout)*time.Second))
	}

	srv := khttp.NewServer(opts...)
	srv.HandlePrefix("/", router)

	return srv
}
