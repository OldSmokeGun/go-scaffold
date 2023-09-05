package controller

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"

	"go-scaffold/internal/app/domain"
	berr "go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/usecase"
)

type PermissionController struct {
	uc   usecase.PermissionUseCaseInterface
	repo repository.PermissionRepositoryInterface
}

func NewPermissionController(
	uc usecase.PermissionUseCaseInterface,
	repo repository.PermissionRepositoryInterface,
) *PermissionController {
	return &PermissionController{
		uc:   uc,
		repo: repo,
	}
}

type PermissionAttr struct {
	Key      string `json:"key"`      // 权限标识
	Name     string `json:"title"`    // 权限
	Desc     string `json:"desc"`     // 权限描述
	ParentID int64  `json:"parentID"` // 父级权限 id
}

func (r PermissionAttr) Validate() error {
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

type PermissionCreateRequest struct {
	PermissionAttr
}

func (r PermissionCreateRequest) toEntity() domain.Permission {
	return domain.Permission{
		Key:      r.Key,
		Name:     r.Name,
		Desc:     r.Desc,
		ParentID: r.ParentID,
	}
}

func (c *PermissionController) Create(ctx context.Context, req PermissionCreateRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	exist, err := c.repo.KeyExist(ctx, req.Key)
	if err != nil {
		return err
	}
	if exist {
		return berr.ErrBadCall.WithMsg("permission key already exist").WithError(errors.New("key already exist"))
	}

	return c.uc.Create(ctx, req.toEntity())
}

type PermissionUpdateRequest struct {
	ID int64 `json:"id"`
	PermissionAttr
}

func (r PermissionUpdateRequest) toEntity() domain.Permission {
	return domain.Permission{
		ID:       r.ID,
		Key:      r.Key,
		Name:     r.Name,
		Desc:     r.Desc,
		ParentID: r.ParentID,
	}
}

func (r PermissionUpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required.Error("id is required")),
		validation.Field(&r.PermissionAttr),
	)
}

func (c *PermissionController) Update(ctx context.Context, req PermissionUpdateRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	permission, err := c.repo.FindOne(ctx, req.ID)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.WithError(err)
	} else if err != nil {
		return err
	}

	if req.ParentID == permission.ID {
		return berr.ErrBadCall.WithMsg("parent cannot be self").WithError(errors.New("parent cannot be self"))
	}

	exist, err := c.repo.KeyExistExcludeID(ctx, req.Key, req.ID)
	if err != nil {
		return err
	}
	if exist {
		return berr.ErrBadCall.WithMsg("permission key already exist").WithError(errors.New("key already exist"))
	}

	return c.uc.Update(ctx, req.toEntity())
}

func (c *PermissionController) Delete(ctx context.Context, id int64) error {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	permission, err := c.repo.FindOne(ctx, id)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.WithError(err)
	} else if err != nil {
		return err
	}

	hasChild, err := c.repo.HasChild(ctx, permission.ID)
	if err != nil {
		return err
	}
	if hasChild {
		return berr.ErrBadCall.WithMsg("permission has child").WithError(errors.New("permission has child"))
	}

	return c.uc.Delete(ctx, *permission)
}

func (c *PermissionController) Detail(ctx context.Context, id int64) (*domain.Permission, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	permission, err := c.uc.Detail(ctx, id)
	if repository.IsNotFound(err) {
		return nil, berr.ErrResourceNotFound.WithError(err)
	} else if err != nil {
		return nil, err
	}

	return permission, nil
}

type PermissionListRequest struct {
	Keyword string
}

func (c *PermissionController) List(ctx context.Context, req PermissionListRequest) ([]*domain.Permission, error) {
	param := usecase.PermissionListParam{Keyword: req.Keyword}
	return c.uc.List(ctx, param)
}
