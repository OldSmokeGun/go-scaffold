package controller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewGreetController,
	NewProducerController,
	NewSystemSessionTokenController,
	NewSystemSessionPermissionController,
	NewSystemSessionController,
	NewSystemUserController,
	NewSystemRoleController,
	NewSystemPermissionController,
	NewProductController,
)
