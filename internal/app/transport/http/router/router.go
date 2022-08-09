package router

import (
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/transport/http/middleware/recover"
	"go-scaffold/internal/app/transport/http/pkg/response"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

var ProviderSet = wire.NewSet(
	New,
	NewAPIGroup,
	NewAPIV1Group,
)

type router struct {
	appConf      *config.App
	httpConf     *config.HTTP
	loggerWriter *rotatelogs.RotateLogs
	zLogger      *zap.Logger
	logger       log.Logger
	apiGroup     *apiGroup
}

// New return *Router
func New(
	appConf *config.App,
	httpConf *config.HTTP,
	loggerWriter *rotatelogs.RotateLogs,
	zLogger *zap.Logger,
	logger log.Logger,
	apiGroup *apiGroup,
) http.Handler {
	if httpConf == nil {
		return nil
	}

	engine := gin.New()

	r := &router{
		appConf:      appConf,
		httpConf:     httpConf,
		loggerWriter: loggerWriter,
		zLogger:      zLogger,
		logger:       logger,
		apiGroup:     apiGroup,
	}
	r.setup()
	r.useMiddlewares(engine)
	r.useRoutes(engine)

	return engine
}

func (r *router) setup() {
	var output io.Writer
	if r.loggerWriter == nil {
		output = os.Stdout
	} else {
		output = io.MultiWriter(os.Stdout, r.loggerWriter)
	}
	gin.DefaultWriter = output
	gin.DefaultErrorWriter = output
	gin.DisableConsoleColor()

	switch r.appConf.Env {
	case config.Local:
		gin.SetMode(gin.DebugMode)
	case config.Test:
		gin.SetMode(gin.TestMode)
	case config.Prod:
		gin.SetMode(gin.ReleaseMode)
	}
}

func (r *router) useMiddlewares(engine *gin.Engine) {
	engine.Use(ginzap.Ginzap(r.zLogger.WithOptions(zap.AddCallerSkip(4)), time.RFC3339, false))
	engine.Use(recover.CustomRecoveryWithZap(r.zLogger.WithOptions(zap.AddCallerSkip(4)), true, func(c *gin.Context, err any) {
		response.Error(c, errors.ServerError())
		c.Abort()
	}))
	engine.Use(otelgin.Middleware(r.appConf.Name))
}

func (r *router) useRoutes(engine *gin.Engine) {
	group := engine.Group("/")
	extAddrSubs := parseExternalAddr(r.httpConf.ExternalAddr)
	if len(extAddrSubs) == 2 {
		group = engine.Group("/" + extAddrSubs[1])
	}

	group.GET("/ping", func(ctx *gin.Context) { ctx.String(http.StatusOK, "pong"); return })

	// register api routing group
	r.apiGroup.setup(group)
	r.apiGroup.useMiddlewares()
	r.apiGroup.useRoutes(engine)
}

// parseExternalAddr parse external address
func parseExternalAddr(addr string) []string {
	return strings.SplitN(addr, "/", 2)
}
