//go:build wireinject
// +build wireinject

package main

import (
	"go-scaffold/internal/app"
	"go-scaffold/internal/app/command"
	appconfig "go-scaffold/internal/app/config"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
)

func initApp(
	*rotatelogs.RotateLogs,
	log.Logger,
	*zap.Logger,
	*appconfig.Config,
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
) (*command.Command, func(), error) {
	panic(wire.Build(
		command.ProviderSet,
		command.New,
	))
}
