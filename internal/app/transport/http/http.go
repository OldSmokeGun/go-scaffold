package http

import (
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/transport/http/handler"
	"go-scaffold/internal/app/transport/http/router"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	handler.ProviderSet,
	router.ProviderSet,
	NewServer,
)

// NewServer 创建 HTTP 服务器
func NewServer(
	logger log.Logger,
	httpConf *config.HTTP,
	router *gin.Engine,
) *khttp.Server {
	if router == nil {
		return nil
	}

	var opts = []khttp.ServerOption{
		khttp.Logger(logger),
	}

	if httpConf.Network != "" {
		opts = append(opts, khttp.Network(httpConf.Network))
	}

	if httpConf.Addr != "" {
		opts = append(opts, khttp.Address(httpConf.Addr))
	}

	if httpConf.Timeout != 0 {
		opts = append(opts, khttp.Timeout(time.Duration(httpConf.Timeout)*time.Second))
	}

	srv := khttp.NewServer(opts...)
	srv.HandlePrefix("/", router)

	return srv
}
