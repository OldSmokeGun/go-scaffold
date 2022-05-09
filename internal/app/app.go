package app

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"go-scaffold/internal/app/component"
	"go-scaffold/internal/app/component/trace"
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/cron"
	"go-scaffold/internal/app/model"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/service"
	"go-scaffold/internal/app/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/gorm"
)

//go:generate swag fmt -g app.go
//go:generate swag init -g app.go -o transport/http/api --parseInternal

// @title                       API 接口文档
// @description                 API 接口文档
// @version                     0.0.0
// @host                        localhost
// @BasePath                    /api
// @schemes                     http https
// @accept                      json
// @accept                      x-www-form-urlencoded
// @produce                     json
// @securityDefinitions.apikey  LoginAuth
// @in                          header
// @name                        Token

var ProviderSet = wire.NewSet(
	config.ProviderSet,
	component.ProviderSet,
	cron.ProviderSet,
	transport.ProviderSet,
	repository.ProviderSet,
	service.ProviderSet,
)

type App struct {
	logger    *log.Helper
	db        *gorm.DB
	trace     *trace.Tracer
	cron      *cron.Cron
	transport *transport.Transport
}

func New(
	logger log.Logger,
	db *gorm.DB,
	trace *trace.Tracer,
	cron *cron.Cron,
	transport *transport.Transport,
) *App {
	return &App{
		logger:    log.NewHelper(logger),
		db:        db,
		trace:     trace,
		cron:      cron,
		transport: transport,
	}
}

// Start 启动应用
func (a *App) Start() (err error) {
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

	// 启动 cron 服务
	if err = a.cron.Start(); err != nil {
		return
	}

	// 启动 transport 服务
	if err = a.transport.Start(); err != nil {
		return
	}

	return nil
}

// Stop 停止应用
func (a *App) Stop(ctx context.Context) (err error) {
	// 关闭 cron 服务
	if err = a.cron.Stop(ctx); err != nil {
		return
	}

	// 关闭 transport 服务
	if err = a.transport.Stop(); err != nil {
		return
	}

	return nil
}
