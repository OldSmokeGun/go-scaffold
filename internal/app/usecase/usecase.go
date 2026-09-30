package usecase

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(SystemSessionUseCaseInterface), new(*SystemSessionUseCase)), NewSystemSessionUseCase),
	wire.NewSet(wire.Bind(new(SystemUserUseCaseInterface), new(*SystemUserUseCase)), NewSystemUserUseCase),
	wire.NewSet(wire.Bind(new(SystemRoleUseCaseInterface), new(*SystemRoleUseCase)), NewSystemRoleUseCase),
	wire.NewSet(wire.Bind(new(SystemPermissionUseCaseInterface), new(*SystemPermissionUseCase)), NewSystemPermissionUseCase),
	wire.NewSet(wire.Bind(new(ProductUseCaseInterface), new(*ProductUseCase)), NewProductUseCase),
)
