package handler

import (
	"go-scaffold/internal/app/transport/http/handler/v1/greet"
	"go-scaffold/internal/app/transport/http/handler/v1/trace"
	"go-scaffold/internal/app/transport/http/handler/v1/user"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(greet.HandlerInterface), new(*greet.Handler)), greet.NewHandler),
	wire.NewSet(wire.Bind(new(trace.HandlerInterface), new(*trace.Handler)), trace.NewHandler),
	wire.NewSet(wire.Bind(new(user.HandlerInterface), new(*user.Handler)), user.NewHandler),
)
