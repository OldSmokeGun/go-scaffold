package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/http/handler"
	"go-scaffold/internal/app/http/router"
	"go.uber.org/zap"
	"net/http"
)

var ProviderSet = wire.NewSet(
	handler.ProviderSet,
	router.New,
	NewServer,
)

type Server struct {
	logger   *zap.Logger
	httpConf *config.HTTP
	router   *gin.Engine

	httpServer *http.Server
}

// NewServer 创建 HTTP 服务器
func NewServer(
	logger *zap.Logger,
	httpConf *config.HTTP,
	router *gin.Engine,
) *Server {
	httpServer := &http.Server{
		Addr:    httpConf.Addr,
		Handler: router,
	}

	return &Server{
		logger:     logger,
		httpConf:   httpConf,
		router:     router,
		httpServer: httpServer,
	}
}

// Start HTTP 服务启动
func (s *Server) Start() (err error) {
	s.logger.Sugar().Infof("http server started on %s", s.httpConf.Addr)

	if err = s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return
	}

	return nil
}

// Stop HTTP 服务关闭
func (s *Server) Stop(ctx context.Context) (err error) {
	if err = s.httpServer.Shutdown(ctx); err != nil {
		return
	}

	s.logger.Sugar().Info("the http server has been shutdown")

	return nil
}
