package web

import (
	"gin-scaffold/global"
	"gin-scaffold/internal/pkg/validator"
	"gin-scaffold/internal/web/config"
	"gin-scaffold/internal/web/router"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

// Setup 是 web 服务的初始化函数
func Setup() (http.Handler, error) {
	err := validator.RegisterValidator([]validator.CustomValidator{
		{"phone", validator.IsPhone},
	})
	if err != nil {
		panic(err)
	}

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

	handler := gin.Default()
	router.Register(handler)

	// ...

	return handler, nil
}

// MustSetup 是 web 服务的初始化函数
func MustSetup() http.Handler {
	handler, err := Setup()
	if err != nil {
		panic(err)
	}

	return handler
}
