package app

import (
	"go-scaffold/internal/app/adapter"
	"go-scaffold/internal/app/controller"
	"go-scaffold/internal/app/pkg"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/usecase"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	adapter.ProviderSet,
	controller.ProviderSet,
	usecase.ProviderSet,
	repository.ProviderSet,
	pkg.ProviderSet,
)
