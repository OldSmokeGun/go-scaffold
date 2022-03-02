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
	cm *config.Config,
	router *gin.Engine,
) *khttp.Server {
	var opts = []khttp.ServerOption{
		khttp.Logger(logger),
	}

	if cm.App.Http.Network != "" {
		opts = append(opts, khttp.Network(cm.App.Http.Network))
	}

	if cm.App.Http.Addr != "" {
		opts = append(opts, khttp.Address(cm.App.Http.Addr))
	}

	if cm.App.Http.Timeout != 0 {
		opts = append(opts, khttp.Timeout(time.Duration(cm.App.Http.Timeout)*time.Second))
	}

	srv := khttp.NewServer(opts...)
	srv.HandlePrefix("/", router)

	return srv
}
