package repository

import (
	"context"

	"github.com/pkg/errors"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/pkg/ent/ent"
	"go-scaffold/internal/app/pkg/ent/ent/product"
)

var _ ProductRepositoryInterface = (*ProductRepository)(nil)

type (
	ProductFindListParam struct {
		Keyword string
	}

	ProductRepositoryInterface interface {
		Filter(ctx context.Context, param ProductFindListParam) ([]*domain.Product, error)
		FindOne(ctx context.Context, id int64) (*domain.Product, error)
		Exist(ctx context.Context, id int64) (bool, error)
		Create(ctx context.Context, e domain.Product) error
		Update(ctx context.Context, e domain.Product) error
		Delete(ctx context.Context, e domain.Product) error
	}
)

type ProductRepository struct {
	client *ent.Client
}

func NewProductRepository(client *ent.Client) *ProductRepository {
	return &ProductRepository{
		client: client,
	}
}

func (r *ProductRepository) Filter(ctx context.Context, param ProductFindListParam) ([]*domain.Product, error) {
	query := r.client.Product.Query()

	if param.Keyword != "" {
		query.Where(
			product.Or(
				product.NameContains(param.Keyword),
				product.DescContains(param.Keyword),
			),
		)
	}

	list, err := query.
		Order(ent.Desc(product.FieldUpdatedAt)).
		All(ctx)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}

	entities := make([]*domain.Product, 0, len(list))
	for _, i := range list {
		entities = append(entities, (&productModel{i}).toEntity())
	}

	return entities, nil
}

func (r *ProductRepository) FindOne(ctx context.Context, id int64) (*domain.Product, error) {
	m, err := r.client.Product.Get(ctx, id)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}
	return (&productModel{m}).toEntity(), nil
}

func (r *ProductRepository) Exist(ctx context.Context, id int64) (bool, error) {
	exist, err := r.client.Product.Query().Where(product.IDEQ(id)).Exist(ctx)
	return exist, errors.WithStack(handleError(err))
}

func (r *ProductRepository) Create(ctx context.Context, e domain.Product) error {
	_, err := r.client.Product.Create().
		SetName(e.Name).
		SetDesc(e.Desc).
		SetPrice(e.Price).
		Save(ctx)
	return errors.WithStack(handleError(err))
}

func (r *ProductRepository) Update(ctx context.Context, e domain.Product) error {
	_, err := r.client.Product.
		UpdateOneID(e.ID).
		SetName(e.Name).
		SetDesc(e.Desc).
		SetPrice(e.Price).
		Save(ctx)
	return errors.WithStack(handleError(err))
}

func (r *ProductRepository) Delete(ctx context.Context, e domain.Product) error {
	return errors.WithStack(r.client.Product.DeleteOneID(e.ID).Exec(ctx))
}

type productModel struct {
	*ent.Product
}

func (m *productModel) toEntity() *domain.Product {
	return &domain.Product{
		ID:    m.ID,
		Name:  m.Name,
		Desc:  m.Desc,
		Price: m.Price,
	}
}
