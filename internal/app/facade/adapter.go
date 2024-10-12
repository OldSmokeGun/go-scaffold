package facade

import (
	"github.com/google/wire"

	"go-scaffold/internal/app/facade/cron"
	"go-scaffold/internal/app/facade/kafka"
	"go-scaffold/internal/app/facade/scripts"
	"go-scaffold/internal/app/facade/server"
)

var ProviderSet = wire.NewSet(
	cron.ProviderSet,
	server.ProviderSet,
	kafka.ProviderSet,
	scripts.ProviderSet,
)
