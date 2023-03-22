package controller

import (
	"context"

	"go-scaffold/internal/app/domain"
	berr "go-scaffold/internal/app/pkg/errors"

	"github.com/pkg/errors"
)

// UserController 用户控制器
type UserController struct {
	uc domain.UserUseCase
}

// NewUserController 构造用户控制器
func NewUserController(
	uc domain.UserUseCase,
) *UserController {
	return &UserController{
		uc: uc,
	}
}

// Create 新增用户
func (c *UserController) Create(ctx context.Context, u domain.User) error {
	if err := u.ValidateWithContext(domain.SetSceneWithContext(ctx, domain.Create)); err != nil {
		return errors.WithStack(berr.ErrValidateError.WithMsg(err.Error()))
	}

	return c.uc.Create(ctx, u)
}

// Update 更新用户
func (c *UserController) Update(ctx context.Context, u domain.User) error {
	if err := u.ValidateWithContext(domain.SetSceneWithContext(ctx, domain.Update)); err != nil {
		return errors.WithStack(berr.ErrValidateError.WithMsg(err.Error()))
	}

	return c.uc.Update(ctx, u)
}

// Delete 删除用户
func (c *UserController) Delete(ctx context.Context, id domain.ID) error {
	if err := id.ValidateWithContext(domain.SetSceneWithContext(ctx, domain.Delete)); err != nil {
		return errors.WithStack(berr.ErrValidateError.WithMsg(err.Error()))
	}

	return c.uc.Delete(ctx, id)
}

// Detail 用户详情
func (c *UserController) Detail(ctx context.Context, id domain.ID) (*domain.User, error) {
	if err := id.ValidateWithContext(domain.SetSceneWithContext(ctx, domain.Detail)); err != nil {
		return nil, errors.WithStack(berr.ErrValidateError.WithMsg(err.Error()))
	}

	return c.uc.Detail(ctx, id)
}

// List 用户列表
func (c *UserController) List(ctx context.Context, req domain.UserListParam) ([]*domain.User, error) {
	return c.uc.List(ctx, req)
}
