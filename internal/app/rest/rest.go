package rest

import (
	"context"
	"fmt"
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/rest/router"
	"net/http"
)

var httpServer *http.Server

// Start HTTP 服务启动
func Start() (err error) {
	// 启动 http 服务
	addr := fmt.Sprintf(
		"%s:%d",
		global.Config().REST.Host,
		global.Config().REST.Port,
	)
	httpServer = &http.Server{
		Addr:    addr,
		Handler: router.New(),
	}

	global.Logger().Sugar().Infof("http server started on %s", addr)

	if err = httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return
	}

	return nil
}

// Stop HTTP 服务关闭
func Stop(ctx context.Context) (err error) {
	if err = httpServer.Shutdown(ctx); err != nil {
		return
	}

	return nil
}
