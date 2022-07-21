package service

import (
	"go-scaffold/internal/app/service/greet"
	"go-scaffold/internal/app/service/user"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(greet.ServiceInterface), new(*greet.Service)), greet.NewService),
	wire.NewSet(wire.Bind(new(user.ServiceInterface), new(*user.Service)), user.NewService),
)
