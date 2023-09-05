package controller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewGreetController,
	NewProducerController,
	NewAccountTokenController,
	NewAccountPermissionController,
	NewAccountController,
	NewUserController,
	NewRoleController,
	NewPermissionController,
	NewProductController,
)
