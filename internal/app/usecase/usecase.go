package usecase

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(AccountUseCaseInterface), new(*AccountUseCase)), NewAccountUseCase),
	wire.NewSet(wire.Bind(new(UserUseCaseInterface), new(*UserUseCase)), NewUserUseCase),
	wire.NewSet(wire.Bind(new(RoleUseCaseInterface), new(*RoleUseCase)), NewRoleUseCase),
	wire.NewSet(wire.Bind(new(PermissionUseCaseInterface), new(*PermissionUseCase)), NewPermissionUseCase),
	wire.NewSet(wire.Bind(new(ProductUseCaseInterface), new(*ProductUseCase)), NewProductUseCase),
)
