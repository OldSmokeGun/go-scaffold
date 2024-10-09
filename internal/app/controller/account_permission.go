package controller

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/pkg/errors"

	"go-scaffold/internal/app/repository"
	berr "go-scaffold/internal/errors"
)

type AccountPermissionController struct {
	roleRepo       repository.RoleRepositoryInterface
	permissionRepo repository.PermissionRepositoryInterface
	enforcer       *casbin.Enforcer
}

func NewAccountPermissionController(
	roleRepo repository.RoleRepositoryInterface,
	permissionRepo repository.PermissionRepositoryInterface,
	enforcer *casbin.Enforcer,
) *AccountPermissionController {
	return &AccountPermissionController{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		enforcer:       enforcer,
	}
}

func (c *AccountPermissionController) ValidatePermission(ctx context.Context, user int64, permissionKey string) (bool, error) {
	permission, err := c.permissionRepo.FindOneByKey(ctx, permissionKey)
	if repository.IsNotFound(err) {
		return false, berr.ErrAccessDenied.WithError(err)
	} else if err != nil {
		return false, err
	}
	result, err := c.enforcer.Enforce(repository.GetPolicyUser(user), fmt.Sprintf("%d", permission.ID))
	if err != nil {
		return false, berr.ErrAccessDenied.WithError(errors.WithStack(err))
	}

	return result, nil
}
