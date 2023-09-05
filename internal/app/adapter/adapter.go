package adapter

import (
	"github.com/google/wire"

	"go-scaffold/internal/app/adapter/cron"
	"go-scaffold/internal/app/adapter/kafka"
	"go-scaffold/internal/app/adapter/server"
)

var ProviderSet = wire.NewSet(
	cron.ProviderSet,
	server.ProviderSet,
	kafka.ProviderSet,
)
