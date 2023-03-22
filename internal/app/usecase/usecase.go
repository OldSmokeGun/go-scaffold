package usecase

import (
	"go-scaffold/internal/app/domain"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(domain.UserUseCase), new(*UserUseCase)), NewUserUseCase),
)
