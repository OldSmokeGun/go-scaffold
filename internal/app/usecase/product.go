package usecase

import (
	"context"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
)

var _ ProductUseCaseInterface = (*ProductUseCase)(nil)

type ProductUseCaseInterface interface {
	Create(ctx context.Context, product domain.Product) error
	Update(ctx context.Context, product domain.Product) error
	Delete(ctx context.Context, product domain.Product) error
	Detail(ctx context.Context, id int64) (*domain.Product, error)
	List(ctx context.Context, param ProductListParam) ([]*domain.Product, error)
}

type ProductUseCase struct {
	repo repository.ProductRepositoryInterface
}

func NewProductUseCase(
	repo repository.ProductRepositoryInterface,
) *ProductUseCase {
	return &ProductUseCase{
		repo: repo,
	}
}

func (c *ProductUseCase) Create(ctx context.Context, product domain.Product) error {
	return c.repo.Create(ctx, product)
}

func (c *ProductUseCase) Update(ctx context.Context, product domain.Product) error {
	return c.repo.Update(ctx, product)
}

func (c *ProductUseCase) Delete(ctx context.Context, product domain.Product) error {
	return c.repo.Delete(ctx, product)
}

func (c *ProductUseCase) Detail(ctx context.Context, id int64) (*domain.Product, error) {
	return c.repo.FindOne(ctx, id)
}

type ProductListParam struct {
	Keyword string
}

func (c *ProductUseCase) List(ctx context.Context, param ProductListParam) ([]*domain.Product, error) {
	return c.repo.Filter(ctx, repository.ProductFindListParam{
		Keyword: param.Keyword,
	})
}
