package usecase

import (
	"context"

	"go-scaffold/internal/app/domain"
	berr "go-scaffold/internal/app/pkg/errors"

	"github.com/pkg/errors"
)

var _ domain.UserUseCase = (*UserUseCase)(nil)

// UserUseCase 用户用例
type UserUseCase struct {
	repo domain.UserRepository
}

// NewUserUseCase 构造用户用例
func NewUserUseCase(
	repo domain.UserRepository,
) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

// Create 新增用户
func (u *UserUseCase) Create(ctx context.Context, user domain.User) error {
	return u.repo.Create(ctx, user)
}

// Update 更新用户
func (u *UserUseCase) Update(ctx context.Context, user domain.User) error {
	ok, err := u.repo.Exist(ctx, user.ID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.WithStack(berr.ErrResourceNotFound)
	}

	return u.repo.Update(ctx, user)
}

// Delete 删除用户
func (u *UserUseCase) Delete(ctx context.Context, id domain.ID) error {
	user, err := u.repo.FindOneByID(ctx, id)
	if err != nil {
		return err
	}

	return u.repo.Delete(ctx, *user)
}

// Detail 用户详情
func (u *UserUseCase) Detail(ctx context.Context, id domain.ID) (*domain.User, error) {
	return u.repo.FindOneByID(ctx, id)
}

// List 用户列表
func (u *UserUseCase) List(ctx context.Context, param domain.UserListParam) ([]*domain.User, error) {
	return u.repo.FindList(ctx, domain.FindUserListParam{
		Keyword: param.Keyword,
	})
}
