package adapter

import (
	"go-scaffold/internal/app/adapter/cron"
	"go-scaffold/internal/app/adapter/server"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	cron.ProviderSet,
	server.ProviderSet,
)
