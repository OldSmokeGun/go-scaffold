package handler

import (
	"go-scaffold/internal/app/command/handler/greet"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(greet.NewHandler)
