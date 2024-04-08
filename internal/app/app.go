package app

import (
	"github.com/google/wire"

	"go-scaffold/internal/app/adapter"
	"go-scaffold/internal/app/controller"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/usecase"
)

var ProviderSet = wire.NewSet(
	adapter.ProviderSet,
	controller.ProviderSet,
	usecase.ProviderSet,
	repository.ProviderSet,
)
