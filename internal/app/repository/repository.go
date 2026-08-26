package repository

import (
	"errors"

	"github.com/google/wire"
	"gorm.io/gorm"

	uerr "go-scaffold/pkg/errors"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(UserRepositoryInterface), new(*UserRepository)), NewUserRepository),
	wire.NewSet(wire.Bind(new(RoleRepositoryInterface), new(*RoleRepository)), NewRoleRepository),
	wire.NewSet(wire.Bind(new(PermissionRepositoryInterface), new(*PermissionRepository)), NewPermissionRepository),
	wire.NewSet(wire.Bind(new(ProductRepositoryInterface), new(*ProductRepository)), NewProductRepository),
)

var ErrRecordNotFound = errors.New("record not found")

func IsNotFound(err error) bool {
	return errors.Is(err, ErrRecordNotFound) || errors.Is(err, gorm.ErrRecordNotFound)
}

func handleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrRecordNotFound
	}

	return uerr.WithStack(err, 4)
}
