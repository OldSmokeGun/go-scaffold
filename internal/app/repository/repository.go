package repository

import (
	"errors"

	"github.com/google/wire"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"

	"go-scaffold/internal/app/pkg/ent/ent"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(UserRepositoryInterface), new(*UserRepository)), NewUserRepository),
	wire.NewSet(wire.Bind(new(RoleRepositoryInterface), new(*RoleRepository)), NewRoleRepository),
	wire.NewSet(wire.Bind(new(PermissionRepositoryInterface), new(*PermissionRepository)), NewPermissionRepository),
	wire.NewSet(wire.Bind(new(ProductRepositoryInterface), new(*ProductRepository)), NewProductRepository),
)

var ErrRecordNotFound = errors.New("record not found")

func IsNotFound(err error) bool {
	return errors.Is(err, ErrRecordNotFound) || errors.Is(err, gorm.ErrRecordNotFound) || ent.IsNotFound(err)
}

// handleError handle ent and gorm error
// masking the internal implementation of the repository layer
func handleError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || ent.IsNotFound(err) {
		return ErrRecordNotFound
	}
	return err
}

// baseModel base model
// automatic update of timestamps, soft delete
type baseModel struct {
	ID        int64                 `gorm:"primaryKey"`
	CreatedAt int64                 `gorm:"NOT NULL"`
	UpdatedAt int64                 `gorm:"NOT NULL;DEFAULT:0"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;NOT NULL;DEFAULT:0"`
}
