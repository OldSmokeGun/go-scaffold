package usecase

import (
	"context"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
)

var _ RoleUseCaseInterface = (*RoleUseCase)(nil)

type RoleUseCaseInterface interface {
	Create(ctx context.Context, product domain.Role) error
	Update(ctx context.Context, product domain.Role) error
	Delete(ctx context.Context, product domain.Role) error
	Detail(ctx context.Context, id int64) (*domain.Role, error)
	List(ctx context.Context, param RoleListParam) ([]*domain.Role, error)
	GrantPermissions(ctx context.Context, role int64, permissions []int64) error
	GetPermissions(ctx context.Context, id int64) ([]*domain.Permission, error)
}

type RoleUseCase struct {
	repo repository.RoleRepositoryInterface
}

func NewRoleUseCase(
	repo repository.RoleRepositoryInterface,
) *RoleUseCase {
	return &RoleUseCase{
		repo: repo,
	}
}

func (c *RoleUseCase) Create(ctx context.Context, product domain.Role) error {
	return c.repo.Create(ctx, product)
}

func (c *RoleUseCase) Update(ctx context.Context, product domain.Role) error {
	return c.repo.Update(ctx, product)
}

func (c *RoleUseCase) Delete(ctx context.Context, product domain.Role) error {
	return c.repo.Delete(ctx, product)
}

func (c *RoleUseCase) Detail(ctx context.Context, id int64) (*domain.Role, error) {
	return c.repo.FindOne(ctx, id)
}

type RoleListParam struct {
	Keyword string
}

func (c *RoleUseCase) List(ctx context.Context, param RoleListParam) ([]*domain.Role, error) {
	return c.repo.Filter(ctx, repository.RoleFindListParam{
		Keyword: param.Keyword,
	})
}

func (c *RoleUseCase) GrantPermissions(ctx context.Context, role int64, permissions []int64) error {
	return c.repo.GrantPermissions(ctx, role, permissions)
}

func (c *RoleUseCase) GetPermissions(ctx context.Context, id int64) ([]*domain.Permission, error) {
	return c.repo.GetPermissions(ctx, id)
}
