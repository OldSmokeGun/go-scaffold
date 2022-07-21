package handler

import (
	greetpb "go-scaffold/internal/app/transport/grpc/api/scaffold/v1/greet"
	userpb "go-scaffold/internal/app/transport/grpc/api/scaffold/v1/user"
	"go-scaffold/internal/app/transport/grpc/handler/v1/greet"
	"go-scaffold/internal/app/transport/grpc/handler/v1/user"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(greetpb.GreetServer), new(*greet.Handler)), greet.NewHandler),
	wire.NewSet(wire.Bind(new(userpb.UserServer), new(*user.Handler)), user.NewHandler),
)
