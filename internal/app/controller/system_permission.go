package controller

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/usecase"
	berr "go-scaffold/internal/errors"
)

type SystemPermissionController struct {
	uc   usecase.SystemPermissionUseCaseInterface
	repo repository.SystemPermissionRepositoryInterface
}

func NewSystemPermissionController(
	uc usecase.SystemPermissionUseCaseInterface,
	repo repository.SystemPermissionRepositoryInterface,
) *SystemPermissionController {
	return &SystemPermissionController{
		uc:   uc,
		repo: repo,
	}
}

type SystemPermissionAttr struct {
	Key      string `json:"key"`      // 权限标识
	Name     string `json:"title"`    // 权限
	Desc     string `json:"desc"`     // 权限描述
	ParentID int64  `json:"parentID"` // 父级权限 id
}

func (r SystemPermissionAttr) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Key,
			validation.Required.Error("key is required"),
			validation.Length(1, 128).Error("key must be 1 ~ 128 characters"),
		),
		validation.Field(&r.Name,
			validation.Length(0, 128).Error("name must be 0 ~ 128 characters"),
		),
		validation.Field(&r.Desc,
			validation.Length(0, 255).Error("description must be 0 ~ 255 characters"),
		),
	)
}

type SystemPermissionCreateRequest struct {
	SystemPermissionAttr
}

func (r SystemPermissionCreateRequest) toEntity() domain.SystemPermission {
	return domain.SystemPermission{
		Key:      r.Key,
		Name:     r.Name,
		Desc:     r.Desc,
		ParentID: r.ParentID,
	}
}

func (c *SystemPermissionController) Create(ctx context.Context, req SystemPermissionCreateRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.Wrap(err)
	}

	exist, err := c.repo.KeyExist(ctx, req.Key)
	if err != nil {
		return err
	}
	if exist {
		return berr.ErrBadCall.Errorf("permission key already exist")
	}

	return c.uc.Create(ctx, req.toEntity())
}

type SystemPermissionUpdateRequest struct {
	ID int64 `json:"id"`
	SystemPermissionAttr
}

func (r SystemPermissionUpdateRequest) toEntity() domain.SystemPermission {
	return domain.SystemPermission{
		ID:       r.ID,
		Key:      r.Key,
		Name:     r.Name,
		Desc:     r.Desc,
		ParentID: r.ParentID,
	}
}

func (r SystemPermissionUpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required.Error("id is required")),
		validation.Field(&r.SystemPermissionAttr),
	)
}

func (c *SystemPermissionController) Update(ctx context.Context, req SystemPermissionUpdateRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.Wrap(err)
	}

	permission, err := c.repo.FindOne(ctx, req.ID)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.Wrap(err)
	} else if err != nil {
		return err
	}

	if req.ParentID == permission.ID {
		return berr.ErrBadCall.Errorf("parent cannot be self")
	}

	exist, err := c.repo.KeyExistExcludeID(ctx, req.Key, req.ID)
	if err != nil {
		return err
	}
	if exist {
		return berr.ErrBadCall.Errorf("permission key already exist")
	}

	return c.uc.Update(ctx, req.toEntity())
}

func (c *SystemPermissionController) Delete(ctx context.Context, id int64) error {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return berr.ErrValidateError.Wrap(err)
	}

	permission, err := c.repo.FindOne(ctx, id)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.Wrap(err)
	} else if err != nil {
		return err
	}

	hasChild, err := c.repo.HasChild(ctx, permission.ID)
	if err != nil {
		return err
	}
	if hasChild {
		return berr.ErrBadCall.Errorf("permission has child")
	}

	return c.uc.Delete(ctx, *permission)
}

func (c *SystemPermissionController) Detail(ctx context.Context, id int64) (*domain.SystemPermission, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.Wrap(err)
	}

	permission, err := c.uc.Detail(ctx, id)
	if repository.IsNotFound(err) {
		return nil, berr.ErrResourceNotFound.Wrap(err)
	} else if err != nil {
		return nil, err
	}

	return permission, nil
}

type SystemPermissionListRequest struct {
	Keyword string
}

func (c *SystemPermissionController) List(ctx context.Context, req SystemPermissionListRequest) ([]*domain.SystemPermission, error) {
	param := usecase.SystemPermissionListParam{Keyword: req.Keyword}
	return c.uc.List(ctx, param)
}
