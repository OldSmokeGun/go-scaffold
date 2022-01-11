package rest

import (
	"context"
	"fmt"
	"gin-scaffold/internal/app/global"
	"gin-scaffold/internal/app/rest/pkg/validatorx"
	"gin-scaffold/internal/app/rest/router"
	"log"
	"net/http"
)

//go:generate swag fmt -g rest.go
//go:generate swag init -g rest.go -o ./api/docs --parseInternal

// @title                       API 接口文档
// @description                 API 接口文档
// @version                     1.0.0
// @host                        localhost:9527
// @BasePath                    /
// @schemes                     http https
// @accept                      json
// @accept                      x-www-form-urlencoded
// @produce                     json
// @securityDefinitions.apikey  LoginAuth
// @in                          header
// @name                        Token

var httpServer *http.Server

// Start HTTP 接口服务启动入口
func Start() (err error) {
	// 注册自定义验证器
	err = validatorx.RegisterCustomValidation([]validatorx.CustomValidator{
		{"phone", validatorx.ValidatePhone},
	})
	if err != nil {
		return
	}

	// 启动 http 服务
	addr := fmt.Sprintf(
		"%s:%d",
		global.Config().App.Host,
		global.Config().App.Port,
	)
	httpServer = &http.Server{
		Addr:    addr,
		Handler: router.New(),
	}

	log.Printf("http server started on %s\n", addr)

	if err = httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return
	}

	return nil
}

// Stop HTTP 接口服务关闭入口
func Stop(ctx context.Context) (err error) {
	if err = httpServer.Shutdown(ctx); err != nil {
		return
	}

	return
}
