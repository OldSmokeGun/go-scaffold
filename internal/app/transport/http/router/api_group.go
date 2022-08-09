package router

import (
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/transport/http/api/docs"
	casbinmd "go-scaffold/internal/app/transport/http/middleware/casbin"
	jwtmd "go-scaffold/internal/app/transport/http/middleware/jwt"
	"go-scaffold/internal/app/transport/http/pkg/response"
	"go-scaffold/internal/app/transport/http/pkg/swagger"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// apiGroup api routing group
type apiGroup struct {
	logger     log.Logger
	appConf    *config.App
	httpConf   *config.HTTP
	jwtConf    *config.JWT
	enforcer   *casbin.Enforcer
	apiV1Group *apiV1Group
	group      *gin.RouterGroup
}

// NewAPIGroup return *apiGroup
func NewAPIGroup(
	logger log.Logger,
	appConf *config.App,
	httpConf *config.HTTP,
	jwtConf *config.JWT,
	enforcer *casbin.Enforcer,
	apiV1Group *apiV1Group,
) *apiGroup {
	return &apiGroup{
		logger:     logger,
		appConf:    appConf,
		httpConf:   httpConf,
		jwtConf:    jwtConf,
		enforcer:   enforcer,
		apiV1Group: apiV1Group,
	}
}

func (g *apiGroup) setup(rg *gin.RouterGroup) {
	g.group = rg.Group("/api")
}

func (g *apiGroup) useMiddlewares() {
	// allowed to cross
	g.group.Use(cors.Default())
	{
		if g.jwtConf != nil {
			if g.jwtConf.Key != "" {
				g.group.Use(jwtmd.New(
					g.jwtConf.Key,
					jwtmd.WithLogger(log.NewHelper(g.logger)),
					jwtmd.WithErrorResponseBody(response.NewBody(int(errors.ServerErrorCode), errors.ServerErrorCode.String(), nil)),
					jwtmd.WithValidateFailedResponseBody(response.NewBody(int(errors.UnauthorizedCode), errors.UnauthorizedCode.String(), nil)),
				).Validate())
			}
		}
		if g.enforcer != nil {
			g.group.Use(casbinmd.New(
				g.enforcer,
				func(ctx *gin.Context) ([]any, error) {
					// TODO
					return nil, nil
				},
				casbinmd.WithLogger(log.NewHelper(g.logger)),
				casbinmd.WithErrorResponseBody(response.NewBody(int(errors.ServerErrorCode), errors.ServerErrorCode.String(), nil)),
				casbinmd.WithValidateFailedResponseBody(response.NewBody(int(errors.UnauthorizedCode), errors.UnauthorizedCode.String(), nil)),
			).Validate())
		}
	}
}

func (g *apiGroup) useRoutes(r *gin.Engine) {
	// register swagger routing
	g.useSwagger(r)

	// register v1 version api routing group
	g.apiV1Group.setup(g.group)
	g.apiV1Group.useRoutes()
}

func (g *apiGroup) useSwagger(r *gin.Engine) {
	// swagger documentation
	if g.appConf.Env == config.Local {
		docs.SwaggerInfo.Host = g.httpConf.Addr
		extAddrSubs := parseExternalAddr(g.httpConf.ExternalAddr)
		if len(extAddrSubs) > 0 {
			docs.SwaggerInfo.Host = extAddrSubs[0]
		}
		docs.SwaggerInfo.BasePath = g.group.BasePath()

		swagger.Register(r, swagger.Config{
			Path: g.group.BasePath() + "/docs",
			Option: func(c *ginSwagger.Config) {
				c.DefaultModelsExpandDepth = -1
			},
		})
	}
}
