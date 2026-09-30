package controller

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"

	"go-scaffold/internal/app/repository"
	berr "go-scaffold/internal/errors"
)

type SystemSessionPermissionController struct {
	roleRepo       repository.SystemRoleRepositoryInterface
	permissionRepo repository.SystemPermissionRepositoryInterface
	enforcer       *casbin.Enforcer
}

func NewSystemSessionPermissionController(
	roleRepo repository.SystemRoleRepositoryInterface,
	permissionRepo repository.SystemPermissionRepositoryInterface,
	enforcer *casbin.Enforcer,
) *SystemSessionPermissionController {
	return &SystemSessionPermissionController{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		enforcer:       enforcer,
	}
}

func (c *SystemSessionPermissionController) ValidatePermission(ctx context.Context, user int64, permissionKey string) (bool, error) {
	permission, err := c.permissionRepo.FindOneByKey(ctx, permissionKey)
	if repository.IsNotFound(err) {
		return false, berr.ErrAccessDenied.Wrap(err)
	} else if err != nil {
		return false, err
	}
	result, err := c.enforcer.Enforce(repository.GetPolicyUser(user), fmt.Sprintf("%d", permission.ID))
	if err != nil {
		return false, berr.ErrAccessDenied.Wrap(err)
	}

	return result, nil
}
