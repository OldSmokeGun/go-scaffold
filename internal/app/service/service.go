package service

import (
	"github.com/google/wire"
	"go-scaffold/internal/app/service/v1/greet"
	"go-scaffold/internal/app/service/v1/user"
)

var ProviderSet = wire.NewSet(
	greet.NewService,
	user.NewService,
)
