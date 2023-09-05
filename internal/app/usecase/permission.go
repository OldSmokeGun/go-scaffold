package usecase

import (
	"context"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
)

var _ PermissionUseCaseInterface = (*PermissionUseCase)(nil)

type PermissionUseCaseInterface interface {
	Create(ctx context.Context, product domain.Permission) error
	Update(ctx context.Context, product domain.Permission) error
	Delete(ctx context.Context, product domain.Permission) error
	Detail(ctx context.Context, id int64) (*domain.Permission, error)
	List(ctx context.Context, param PermissionListParam) ([]*domain.Permission, error)
}

type PermissionUseCase struct {
	repo repository.PermissionRepositoryInterface
}

func NewPermissionUseCase(
	repo repository.PermissionRepositoryInterface,
) *PermissionUseCase {
	return &PermissionUseCase{
		repo: repo,
	}
}

func (c *PermissionUseCase) Create(ctx context.Context, product domain.Permission) error {
	return c.repo.Create(ctx, product)
}

func (c *PermissionUseCase) Update(ctx context.Context, product domain.Permission) error {
	return c.repo.Update(ctx, product)
}

func (c *PermissionUseCase) Delete(ctx context.Context, product domain.Permission) error {
	return c.repo.Delete(ctx, product)
}

func (c *PermissionUseCase) Detail(ctx context.Context, id int64) (*domain.Permission, error) {
	return c.repo.FindOne(ctx, id)
}

type PermissionListParam struct {
	Keyword string
}

func (c *PermissionUseCase) List(ctx context.Context, param PermissionListParam) ([]*domain.Permission, error) {
	return c.repo.Filter(ctx, repository.PermissionFindListParam{
		Keyword: param.Keyword,
	})
}
