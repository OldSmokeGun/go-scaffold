package web

import (
	"gin-scaffold/internal/web/config"
	"gin-scaffold/internal/web/global"
	"gin-scaffold/internal/web/pkg/validatorx"
	"gin-scaffold/internal/web/routes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

//go:generate swag fmt -g web.go
//go:generate swag init -g web.go -o ./docs --parseInternal

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

// Setup 是 web 服务的初始化函数
func Setup() (http.Handler, error) {
	// 初始化 router
	gin.DefaultWriter = io.MultiWriter(global.LoggerOutput(), os.Stdout)

	switch global.Config().Env {
	case config.Local:
		gin.SetMode(gin.DebugMode)
	case config.Test:
		gin.SetMode(gin.TestMode)
	case config.Prod:
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	routes.Register(router)

	// 注册自定义验证器
	err := validatorx.RegisterCustomValidation([]validatorx.CustomValidator{
		{"phone", validatorx.ValidatePhone},
	})
	if err != nil {
		panic(err)
	}

	// ...

	return router, nil
}

// MustSetup 是 web 服务的初始化函数
func MustSetup() http.Handler {
	handler, err := Setup()
	if err != nil {
		panic(err)
	}

	return handler
}
