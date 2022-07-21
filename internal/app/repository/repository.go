package repository

import (
	"go-scaffold/internal/app/repository/user"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(user.RepositoryInterface), new(*user.Repository)), user.NewRepository),
)
