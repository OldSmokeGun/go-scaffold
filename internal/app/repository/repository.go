package repository

import (
	"github.com/google/wire"
	"go-scaffold/internal/app/repository/user"
)

var ProviderSet = wire.NewSet(user.New)
