package router

import (
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/google/wire"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/transport/http/api"
	"go-scaffold/internal/app/transport/http/handler/v1/greet"
	"go-scaffold/internal/app/transport/http/handler/v1/trace"
	"go-scaffold/internal/app/transport/http/handler/v1/user"
	"go-scaffold/internal/app/transport/http/middleware/recover"
	"go-scaffold/internal/app/transport/http/pkg/response"
	"go-scaffold/internal/app/transport/http/pkg/swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"io"
	"os"
)

var ProviderSet = wire.NewSet(
	New,
)

// New 返回 gin 路由对象
func New(
	loggerWriter *rotatelogs.RotateLogs,
	zLogger *zap.Logger,
	appConf *config.App,
	greetHandler greet.HandlerInterface,
	traceHandler trace.HandlerInterface,
	userHandler user.HandlerInterface,
) *gin.Engine {
	if appConf.Http == nil {
		return nil
	}

	var output io.Writer
	if loggerWriter == nil {
		output = os.Stdout
	} else {
		output = io.MultiWriter(os.Stdout, loggerWriter)
	}
	gin.DefaultWriter = output
	gin.DefaultErrorWriter = output
	gin.DisableConsoleColor()

	switch appConf.Env {
	case config.Local:
		gin.SetMode(gin.DebugMode)
	case config.Test:
		gin.SetMode(gin.TestMode)
	case config.Prod:
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(ginzap.Ginzap(zLogger.WithOptions(zap.AddCallerSkip(4)), time.RFC3339, false))
	router.Use(recover.CustomRecoveryWithZap(zLogger.WithOptions(zap.AddCallerSkip(4)), true, func(c *gin.Context, err interface{}) {
		response.Error(c, errors.ServerError())
		c.Abort()
	}))
	router.Use(otelgin.Middleware(appConf.Name))

	rg := router.Group("/")
	extAddrSubs := strings.SplitN(appConf.Http.ExternalAddr, "/", 2)
	if len(extAddrSubs) == 2 {
		rg = router.Group("/" + extAddrSubs[1])
	}

	rg.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
		return
	})
	// 注册 api 路由组
	apiGroup := rg.Group("/api")
	{
		apiGroup.Use(
			cors.Default(), // 允许跨越
			// jwt.Validate(
			// 	appConf.App.Jwt.Key,
			// 	jwt.WithErrorResponseBody(responsex.NewServerErrorBody()),
			// 	jwt.WithValidateErrorResponseBody(responsex.NewUnauthorizedBody()),
			// 	jwt.WithLogger(log.NewHelper(logger)),
			// ), // jwt 认证
		)

		// swagger 配置
		if appConf.Env == config.Local {
			api.SwaggerInfo.Host = appConf.Http.Addr
			if len(extAddrSubs) > 0 {
				api.SwaggerInfo.Host = extAddrSubs[0]
			}
			api.SwaggerInfo.BasePath = apiGroup.BasePath()

			swagger.Setup(router, swagger.Config{
				Path: apiGroup.BasePath() + "/docs",
				Option: func(c *ginSwagger.Config) {
					c.DefaultModelsExpandDepth = -1
				},
			})
		}

		apiV1Group := apiGroup.Group("/v1")
		{
			// TODO 编写路由

			apiV1Group.GET("/greet", greetHandler.Hello)
			apiV1Group.GET("/trace", traceHandler.Example)

			apiV1Group.GET("/users", userHandler.List)
			apiV1Group.GET("/user/:id", userHandler.Detail)
			apiV1Group.POST("/user", userHandler.Create)
			apiV1Group.PUT("/user/:id", userHandler.Update)
			apiV1Group.DELETE("/user/:id", userHandler.Delete)
		}
	}

	return router
}
