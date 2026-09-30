package controller

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/samber/lo"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/usecase"
	berr "go-scaffold/internal/errors"
)

type SystemRoleController struct {
	uc             usecase.SystemRoleUseCaseInterface
	roleRepo       repository.SystemRoleRepositoryInterface
	permissionRepo repository.SystemPermissionRepositoryInterface
}

func NewSystemRoleController(
	uc usecase.SystemRoleUseCaseInterface,
	roleRepo repository.SystemRoleRepositoryInterface,
	permissionRepo repository.SystemPermissionRepositoryInterface,
) *SystemRoleController {
	return &SystemRoleController{
		uc:             uc,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

type SystemRoleAttr struct {
	Name string `json:"name"`
}

func (r SystemRoleAttr) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required.Error("name is required"),
			validation.Length(1, 32).Error("name must be 1 ~ 32 characters"),
		),
	)
}

type SystemRoleCreateRequest struct {
	SystemRoleAttr
}

func (r SystemRoleCreateRequest) toEntity() domain.SystemRole {
	return domain.SystemRole{
		Name: r.Name,
	}
}

func (c *SystemRoleController) Create(ctx context.Context, req SystemRoleCreateRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.Wrap(err)
	}

	exist, err := c.roleRepo.NameExist(ctx, req.Name)
	if err != nil {
		return err
	}
	if exist {
		return berr.ErrBadCall.Errorf("role name already exist")
	}

	return c.uc.Create(ctx, req.toEntity())
}

type SystemRoleUpdateRequest struct {
	ID int64 `json:"id"`
	SystemRoleAttr
}

func (r SystemRoleUpdateRequest) toEntity() domain.SystemRole {
	return domain.SystemRole{
		ID:   r.ID,
		Name: r.Name,
	}
}

func (r SystemRoleUpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required.Error("id is required")),
		validation.Field(&r.SystemRoleAttr),
	)
}

func (c *SystemRoleController) Update(ctx context.Context, req SystemRoleUpdateRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.Wrap(err)
	}

	_, err := c.roleRepo.FindOne(ctx, req.ID)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.Wrap(err)
	} else if err != nil {
		return err
	}

	exist, err := c.roleRepo.NameExistExcludeID(ctx, req.Name, req.ID)
	if err != nil {
		return err
	}
	if exist {
		return berr.ErrBadCall.Errorf("role name already exist")
	}

	return c.uc.Update(ctx, req.toEntity())
}

func (c *SystemRoleController) Delete(ctx context.Context, id int64) error {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return berr.ErrValidateError.Wrap(err)
	}

	role, err := c.roleRepo.FindOne(ctx, id)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.Wrap(err)
	} else if err != nil {
		return err
	}

	return c.uc.Delete(ctx, *role)
}

func (c *SystemRoleController) Detail(ctx context.Context, id int64) (*domain.SystemRole, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.Wrap(err)
	}

	role, err := c.uc.Detail(ctx, id)
	if repository.IsNotFound(err) {
		return nil, berr.ErrResourceNotFound.Wrap(err)
	} else if err != nil {
		return nil, err
	}

	return role, nil
}

type SystemRoleListRequest struct {
	Keyword string
}

func (c *SystemRoleController) List(ctx context.Context, req SystemRoleListRequest) ([]*domain.SystemRole, error) {
	param := usecase.SystemRoleListParam{Keyword: req.Keyword}
	return c.uc.List(ctx, param)
}

type SystemRoleGrantPermissionsRequest struct {
	Role        int64
	Permissions []int64
}

func (r SystemRoleGrantPermissionsRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Role, validation.Required.Error("role is required")),
		validation.Field(&r.Permissions, validation.Required.Error("no permissions that will be granted")),
	)
}

func (c *SystemRoleController) GrantPermissions(ctx context.Context, req SystemRoleGrantPermissionsRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.Wrap(err)
	}

	if err := c.validatePermissionsExist(ctx, req.Permissions); err != nil {
		return err
	}

	return c.uc.GrantPermissions(ctx, req.Role, req.Permissions)
}

func (c *SystemRoleController) GetPermissions(ctx context.Context, id int64) ([]*domain.SystemPermission, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.Wrap(err)
	}

	return c.uc.GetPermissions(ctx, id)
}

func (c *SystemRoleController) validatePermissionsExist(ctx context.Context, permissions []int64) error {
	list, err := c.permissionRepo.FindList(ctx, permissions)
	if err != nil {
		return err
	}
	permissionList := lo.Map(list, func(item *domain.SystemPermission, index int) int64 {
		return item.ID
	})

	diffs, _ := lo.Difference(permissions, permissionList)
	if len(diffs) > 0 {
		return berr.ErrBadCall.Errorf("permissions %v not exist", diffs)
	}
	return nil
}
