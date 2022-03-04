package handler

import (
	"github.com/google/wire"
	"go-scaffold/internal/app/transport/http/handler/v1/greet"
	"go-scaffold/internal/app/transport/http/handler/v1/trace"
	"go-scaffold/internal/app/transport/http/handler/v1/user"
)

var ProviderSet = wire.NewSet(
	greet.New,
	trace.New,
	user.New,
)
