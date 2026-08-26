package repository

import (
	"context"

	"gorm.io/gorm"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository/queries"
	"go-scaffold/internal/app/repository/schema"
	"go-scaffold/internal/app/repository/tools"
	igorm "go-scaffold/internal/pkg/gorm"
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
	db *igorm.DefaultDB
}

func NewProductRepository(db *igorm.DefaultDB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Filter(ctx context.Context, param ProductFindListParam) ([]*domain.Product, error) {
	q := gorm.G[schema.Product](r.db).
		Order(queries.Product.UpdatedAt.Desc())

	if param.Keyword != "" {
		kw := tools.BuildLikeContains(param.Keyword)
		q = q.
			Where(queries.Product.Name.Like(kw)).
			Or(queries.Product.Desc.Like(kw))
	}

	list, err := q.Find(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	entities := make([]*domain.Product, 0, len(list))
	for i := range list {
		entities = append(entities, list[i].ToEntity())
	}

	return entities, nil
}

func (r *ProductRepository) FindOne(ctx context.Context, id int64) (*domain.Product, error) {
	m, err := queries.Query[schema.Product](r.db).FindOne(ctx, id)
	if err == nil && m.ID == 0 {
		err = gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, handleError(err)
	}

	return m.ToEntity(), nil
}

func (r *ProductRepository) Exist(ctx context.Context, id int64) (bool, error) {
	n, err := queries.Query[schema.Product](r.db).Exist(ctx, id)

	return n > 0, handleError(err)
}

func (r *ProductRepository) Create(ctx context.Context, e domain.Product) error {
	m := schema.Product{
		Name:  e.Name,
		Desc:  e.Desc,
		Price: e.Price,
	}

	err := gorm.G[schema.Product](r.db).
		Create(ctx, &m)

	return handleError(err)
}

func (r *ProductRepository) Update(ctx context.Context, e domain.Product) error {
	_, err := gorm.G[schema.Product](r.db).
		Where(queries.Product.ID.Eq(e.ID)).
		Set(
			queries.Product.Name.Set(e.Name),
			queries.Product.Desc.Set(e.Desc),
			queries.Product.Price.Set(e.Price),
		).
		Update(ctx)

	return handleError(err)
}

func (r *ProductRepository) Delete(ctx context.Context, e domain.Product) error {
	_, err := gorm.G[schema.Product](r.db).
		Where(queries.Product.ID.Eq(e.ID)).
		Delete(ctx)

	return handleError(err)
}
