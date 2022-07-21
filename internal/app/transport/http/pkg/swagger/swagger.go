package swagger

import (
	_ "go-scaffold/internal/app/transport/http/api/docs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Config swagger 文档配置
type Config struct {
	Path   string
	Option func(c *ginSwagger.Config)
}

// Setup 初始化 swagger 文档
func Setup(router *gin.Engine, conf Config) {
	g := router.Group(conf.Path).Use(func(context *gin.Context) {
		if strings.HasSuffix(context.Request.URL.Path, "/") {
			context.Request.URL.Path = strings.TrimSuffix(context.Request.URL.Path, "/")
			router.HandleContext(context)
			context.Abort()
			return
		}
	})

	g.GET("", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, conf.Path+"/index.html")
	})

	g.GET(
		"/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			conf.Option,
		),
	)
}
