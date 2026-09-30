package usecase

import (
	"context"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
)

var _ SystemRoleUseCaseInterface = (*SystemRoleUseCase)(nil)

type SystemRoleUseCaseInterface interface {
	Create(ctx context.Context, product domain.SystemRole) error
	Update(ctx context.Context, product domain.SystemRole) error
	Delete(ctx context.Context, product domain.SystemRole) error
	Detail(ctx context.Context, id int64) (*domain.SystemRole, error)
	List(ctx context.Context, param SystemRoleListParam) ([]*domain.SystemRole, error)
	GrantPermissions(ctx context.Context, role int64, permissions []int64) error
	GetPermissions(ctx context.Context, id int64) ([]*domain.SystemPermission, error)
}

type SystemRoleUseCase struct {
	repo repository.SystemRoleRepositoryInterface
}

func NewSystemRoleUseCase(
	repo repository.SystemRoleRepositoryInterface,
) *SystemRoleUseCase {
	return &SystemRoleUseCase{
		repo: repo,
	}
}

func (c *SystemRoleUseCase) Create(ctx context.Context, product domain.SystemRole) error {
	return c.repo.Create(ctx, product)
}

func (c *SystemRoleUseCase) Update(ctx context.Context, product domain.SystemRole) error {
	return c.repo.Update(ctx, product)
}

func (c *SystemRoleUseCase) Delete(ctx context.Context, product domain.SystemRole) error {
	return c.repo.Delete(ctx, product)
}

func (c *SystemRoleUseCase) Detail(ctx context.Context, id int64) (*domain.SystemRole, error) {
	return c.repo.FindOne(ctx, id)
}

type SystemRoleListParam struct {
	Keyword string
}

func (c *SystemRoleUseCase) List(ctx context.Context, param SystemRoleListParam) ([]*domain.SystemRole, error) {
	return c.repo.Filter(ctx, repository.SystemRoleFindListParam{
		Keyword: param.Keyword,
	})
}

func (c *SystemRoleUseCase) GrantPermissions(ctx context.Context, role int64, permissions []int64) error {
	return c.repo.GrantPermissions(ctx, role, permissions)
}

func (c *SystemRoleUseCase) GetPermissions(ctx context.Context, id int64) ([]*domain.SystemPermission, error) {
	return c.repo.GetPermissions(ctx, id)
}
