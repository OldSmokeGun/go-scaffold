package service

import (
	"github.com/google/wire"
	greetpb "go-scaffold/internal/app/api/v1/greet"
	userpb "go-scaffold/internal/app/api/v1/user"
	"go-scaffold/internal/app/service/v1/greet"
	"go-scaffold/internal/app/service/v1/user"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(greetpb.GreetServer), new(*greet.Service)), greet.NewService),
	wire.NewSet(wire.Bind(new(userpb.UserServer), new(*user.Service)), user.NewService),
)
