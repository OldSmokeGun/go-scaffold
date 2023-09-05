package controller

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"go-scaffold/internal/app/domain"
	berr "go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/usecase"
)

type RoleController struct {
	uc             usecase.RoleUseCaseInterface
	roleRepo       repository.RoleRepositoryInterface
	permissionRepo repository.PermissionRepositoryInterface
}

func NewRoleController(
	uc usecase.RoleUseCaseInterface,
	roleRepo repository.RoleRepositoryInterface,
	permissionRepo repository.PermissionRepositoryInterface,
) *RoleController {
	return &RoleController{
		uc:             uc,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

type RoleAttr struct {
	Name string `json:"name"`
}

func (r RoleAttr) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required.Error("name is required"),
			validation.Length(1, 32).Error("name must be 1 ~ 32 characters"),
		),
	)
}

type RoleCreateRequest struct {
	RoleAttr
}

func (r RoleCreateRequest) toEntity() domain.Role {
	return domain.Role{
		Name: r.Name,
	}
}

func (c *RoleController) Create(ctx context.Context, req RoleCreateRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	exist, err := c.roleRepo.NameExist(ctx, req.Name)
	if err != nil {
		return err
	}
	if exist {
		return berr.ErrBadCall.WithMsg("role name already exist").WithError(errors.New("name already exist"))
	}

	return c.uc.Create(ctx, req.toEntity())
}

type RoleUpdateRequest struct {
	ID int64 `json:"id"`
	RoleAttr
}

func (r RoleUpdateRequest) toEntity() domain.Role {
	return domain.Role{
		ID:   r.ID,
		Name: r.Name,
	}
}

func (r RoleUpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required.Error("id is required")),
		validation.Field(&r.RoleAttr),
	)
}

func (c *RoleController) Update(ctx context.Context, req RoleUpdateRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	_, err := c.roleRepo.FindOne(ctx, req.ID)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.WithError(err)
	} else if err != nil {
		return err
	}

	exist, err := c.roleRepo.NameExistExcludeID(ctx, req.Name, req.ID)
	if err != nil {
		return err
	}
	if exist {
		return berr.ErrBadCall.WithMsg("role name already exist").WithError(errors.New("name already exist"))
	}

	return c.uc.Update(ctx, req.toEntity())
}

func (c *RoleController) Delete(ctx context.Context, id int64) error {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	role, err := c.roleRepo.FindOne(ctx, id)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.WithError(err)
	} else if err != nil {
		return err
	}

	return c.uc.Delete(ctx, *role)
}

func (c *RoleController) Detail(ctx context.Context, id int64) (*domain.Role, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	role, err := c.uc.Detail(ctx, id)
	if repository.IsNotFound(err) {
		return nil, berr.ErrResourceNotFound.WithError(err)
	} else if err != nil {
		return nil, err
	}

	return role, nil
}

type RoleListRequest struct {
	Keyword string
}

func (c *RoleController) List(ctx context.Context, req RoleListRequest) ([]*domain.Role, error) {
	param := usecase.RoleListParam{Keyword: req.Keyword}
	return c.uc.List(ctx, param)
}

type RoleGrantPermissionsRequest struct {
	Role        int64
	Permissions []int64
}

func (r RoleGrantPermissionsRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Role, validation.Required.Error("role is required")),
		validation.Field(&r.Permissions, validation.Required.Error("no permissions that will be granted")),
	)
}

func (c *RoleController) GrantPermissions(ctx context.Context, req RoleGrantPermissionsRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	if err := c.validatePermissionsExist(ctx, req.Permissions); err != nil {
		return err
	}

	return c.uc.GrantPermissions(ctx, req.Role, req.Permissions)
}

func (c *RoleController) GetPermissions(ctx context.Context, id int64) ([]*domain.Permission, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	return c.uc.GetPermissions(ctx, id)
}

func (c *RoleController) validatePermissionsExist(ctx context.Context, permissions []int64) error {
	list, err := c.permissionRepo.FindList(ctx, permissions)
	if err != nil {
		return err
	}
	permissionList := lo.Map(list, func(item *domain.Permission, index int) int64 {
		return item.ID
	})

	diffs, _ := lo.Difference(permissions, permissionList)
	if len(diffs) > 0 {
		return berr.ErrBadCall.WithMsg(fmt.Sprintf("permissions %v not exist", diffs)).WithError(errors.New("permission not exist"))
	}
	return nil
}
