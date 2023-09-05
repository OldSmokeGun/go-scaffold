package usecase

import (
	"context"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
)

var _ UserUseCaseInterface = (*UserUseCase)(nil)

type UserUseCaseInterface interface {
	Create(ctx context.Context, user domain.User) (*domain.User, error)
	Update(ctx context.Context, user domain.User) (*domain.User, error)
	Delete(ctx context.Context, user domain.User) error
	Detail(ctx context.Context, id int64) (*domain.User, error)
	List(ctx context.Context, param UserListParam) ([]*domain.User, error)
	AssignRoles(ctx context.Context, user int64, roles []int64) error
	GetRoles(ctx context.Context, user int64) ([]*domain.Role, error)
	GetPermissions(ctx context.Context, user int64) ([]*domain.Permission, error)
}

type UserUseCase struct {
	repo repository.UserRepositoryInterface
}

func NewUserUseCase(
	repo repository.UserRepositoryInterface,
) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (c *UserUseCase) Create(ctx context.Context, user domain.User) (*domain.User, error) {
	return c.repo.Create(ctx, user)
}

func (c *UserUseCase) Update(ctx context.Context, user domain.User) (*domain.User, error) {
	return c.repo.Update(ctx, user)
}

func (c *UserUseCase) Delete(ctx context.Context, user domain.User) error {
	return c.repo.Delete(ctx, user)
}

func (c *UserUseCase) Detail(ctx context.Context, id int64) (*domain.User, error) {
	return c.repo.FindOne(ctx, id)
}

type UserListParam struct {
	Keyword string
}

func (c *UserUseCase) List(ctx context.Context, param UserListParam) ([]*domain.User, error) {
	return c.repo.Filter(ctx, repository.UserFindListParam{
		Keyword: param.Keyword,
	})
}

func (c *UserUseCase) AssignRoles(ctx context.Context, user int64, roles []int64) error {
	return c.repo.AssignRoles(ctx, user, roles)
}

func (c *UserUseCase) GetRoles(ctx context.Context, user int64) ([]*domain.Role, error) {
	return c.repo.GetRoles(ctx, user)
}

func (c *UserUseCase) GetPermissions(ctx context.Context, user int64) ([]*domain.Permission, error) {
	return c.repo.GetPermissions(ctx, user)
}
