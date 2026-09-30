package usecase

import (
	"context"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
)

var _ SystemPermissionUseCaseInterface = (*SystemPermissionUseCase)(nil)

type SystemPermissionUseCaseInterface interface {
	Create(ctx context.Context, product domain.SystemPermission) error
	Update(ctx context.Context, product domain.SystemPermission) error
	Delete(ctx context.Context, product domain.SystemPermission) error
	Detail(ctx context.Context, id int64) (*domain.SystemPermission, error)
	List(ctx context.Context, param SystemPermissionListParam) ([]*domain.SystemPermission, error)
}

type SystemPermissionUseCase struct {
	repo repository.SystemPermissionRepositoryInterface
}

func NewSystemPermissionUseCase(
	repo repository.SystemPermissionRepositoryInterface,
) *SystemPermissionUseCase {
	return &SystemPermissionUseCase{
		repo: repo,
	}
}

func (c *SystemPermissionUseCase) Create(ctx context.Context, product domain.SystemPermission) error {
	return c.repo.Create(ctx, product)
}

func (c *SystemPermissionUseCase) Update(ctx context.Context, product domain.SystemPermission) error {
	return c.repo.Update(ctx, product)
}

func (c *SystemPermissionUseCase) Delete(ctx context.Context, product domain.SystemPermission) error {
	return c.repo.Delete(ctx, product)
}

func (c *SystemPermissionUseCase) Detail(ctx context.Context, id int64) (*domain.SystemPermission, error) {
	return c.repo.FindOne(ctx, id)
}

type SystemPermissionListParam struct {
	Keyword string
}

func (c *SystemPermissionUseCase) List(ctx context.Context, param SystemPermissionListParam) ([]*domain.SystemPermission, error) {
	return c.repo.Filter(ctx, repository.SystemPermissionFindListParam{
		Keyword: param.Keyword,
	})
}
