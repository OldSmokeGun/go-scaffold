//go:build wireinject
// +build wireinject

package command

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/wire"

	"go-scaffold/internal/app"
	"go-scaffold/internal/app/adapter/cron"
	"go-scaffold/internal/app/adapter/kafka"
	"go-scaffold/internal/app/adapter/server"
	"go-scaffold/internal/app/pkg"
	"go-scaffold/internal/config"
	"go-scaffold/pkg/trace"
)

func initServer(
	context.Context,
	config.AppName,
	config.Env,
	*slog.Logger,
	*trace.Trace,
) (*server.Server, func(), error) {
	panic(wire.Build(
		config.ProviderSet,
		app.ProviderSet,
	))
}

func initCron(
	context.Context,
	config.AppName,
	config.Env,
	*slog.Logger,
) (*cron.Cron, func(), error) {
	panic(wire.Build(
		config.ProviderSet,
		app.ProviderSet,
	))
}

func initKafka(
	context.Context,
	config.AppName,
	config.Env,
	*slog.Logger,
) (*kafka.Kafka, func(), error) {
	panic(wire.Build(
		config.ProviderSet,
		app.ProviderSet,
	))
}

func initDB(
	context.Context,
	config.DatabaseConn,
	*slog.Logger,
) (*sql.DB, func(), error) {
	panic(wire.Build(
		pkg.ProviderSet,
	))
}
