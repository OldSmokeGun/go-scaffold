//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go-scaffold/internal/app"
	"go-scaffold/internal/app/command"
	"go-scaffold/internal/app/component/data"
	"go-scaffold/internal/app/component/discovery"
	"go-scaffold/internal/app/component/orm"
	"go-scaffold/internal/app/component/redis"
	"go-scaffold/internal/app/component/trace"
	appconfig "go-scaffold/internal/app/config"
	"go.uber.org/zap"
)

func initApp(
	*rotatelogs.RotateLogs,
	log.Logger,
	*zap.Logger,
	*appconfig.Config,
	*orm.Config,
	*data.Config,
	*redis.Config,
	*trace.Config,
	*discovery.Config,
) (*app.App, func(), error) {
	panic(wire.Build(
		app.ProviderSet,
		app.New,
	))
}

func initCommand(
	*rotatelogs.RotateLogs,
	log.Logger,
	*zap.Logger,
	*appconfig.Config,
	*orm.Config,
	*data.Config,
	*redis.Config,
	*trace.Config,
	*discovery.Config,
) (*command.Command, func(), error) {
	panic(wire.Build(
		command.ProviderSet,
		command.New,
	))
}
