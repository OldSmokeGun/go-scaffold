package repository

import (
	"errors"

	"github.com/google/wire"
	"gorm.io/gorm"

	uerr "go-scaffold/pkg/errors"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(SystemUserRepositoryInterface), new(*SystemUserRepository)), NewSystemUserRepository),
	wire.NewSet(wire.Bind(new(SystemRoleRepositoryInterface), new(*SystemRoleRepository)), NewSystemRoleRepository),
	wire.NewSet(wire.Bind(new(SystemPermissionRepositoryInterface), new(*SystemPermissionRepository)), NewSystemPermissionRepository),
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
