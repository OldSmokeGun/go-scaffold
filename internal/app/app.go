package app

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
	"go-scaffold/internal/app/component"
	"go-scaffold/internal/app/component/trace"
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/cron"
	"go-scaffold/internal/app/http"
	"go-scaffold/internal/app/model"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:generate swag fmt -g app.go
//go:generate swag init -g app.go -o http/api/docs --parseInternal

// @title                       API 接口文档
// @description                 API 接口文档
// @version                     0.0.0
// @host                        localhost
// @BasePath                    /api
// @schemes                     http https
// @accept                      json
// @accept                      x-www-form-urlencoded
// @produce                     json
// @securityDefinitions.apikey  Authorization
// @in                          header
// @name                        Token

var ProviderSet = wire.NewSet(
	config.ProviderSet,
	component.ProviderSet,
	cron.ProviderSet,
	http.ProviderSet,
	repository.ProviderSet,
	service.ProviderSet,
)

type App struct {
	logger   *zap.Logger
	db       *gorm.DB
	trace    *trace.Tracer
	cron     *cron.Cron
	http     *http.Server
	enforcer *casbin.Enforcer
}

func New(
	logger *zap.Logger,
	db *gorm.DB,
	trace *trace.Tracer,
	cron *cron.Cron,
	http *http.Server,
	enforcer *casbin.Enforcer,
) *App {
	return &App{
		logger:   logger,
		db:       db,
		trace:    trace,
		cron:     cron,
		http:     http,
		enforcer: enforcer,
	}
}

// Start 启动 app
func (a *App) Start(cancel context.CancelFunc) (err error) {
	// 设置 tracer
	if a.trace != nil {
		otel.SetTracerProvider(a.trace.TracerProvider())
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))
	}

	// 数据迁移
	if a.db != nil {
		if err = model.Migrate(a.db); err != nil {
			return
		}

		a.logger.Info("database migration completed")
	}

	if a.enforcer != nil {
		if err = a.enforcer.LoadPolicy(); err != nil {
			return
		}

		a.logger.Info("casbin policy loaded")
	}

	// 启动 cron 服务
	if err = a.cron.Start(); err != nil {
		return
	}

	// 启动 HTTP 服务
	if err = a.http.Start(); err != nil {
		a.logger.Sugar().Error(err)
		cancel()
		return
	}

	return nil
}

// Stop 关闭 app
func (a *App) Stop(ctx context.Context) (err error) {
	// 关闭 cron 服务
	if err = a.cron.Stop(ctx); err != nil {
		return
	}

	// 关闭 HTTP 服务
	if err = a.http.Stop(ctx); err != nil {
		return
	}

	return nil
}
