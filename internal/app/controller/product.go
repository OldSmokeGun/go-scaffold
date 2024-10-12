package controller

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/usecase"
	berr "go-scaffold/internal/errors"
)

type ProductController struct {
	uc   usecase.ProductUseCaseInterface
	repo repository.ProductRepositoryInterface
}

func NewProductController(
	uc usecase.ProductUseCaseInterface,
	repo repository.ProductRepositoryInterface,
) *ProductController {
	return &ProductController{
		uc:   uc,
		repo: repo,
	}
}

type ProductAttr struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int    `json:"price"`
}

func (r ProductAttr) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required.Error("name is required"),
			validation.Length(8, 128).Error("name must be 8 ~ 128 characters"),
		),
		validation.Field(&r.Desc,
			validation.Required.Error("description is required"),
			validation.Length(0, 255).Error("description must be 0 ~ 255 characters"),
		),
		validation.Field(&r.Price, validation.Required.Error("price is required")),
	)
}

type ProductCreateRequest struct {
	ProductAttr
}

func (r ProductCreateRequest) toEntity() domain.Product {
	return domain.Product{
		Name:  r.Name,
		Desc:  r.Desc,
		Price: r.Price,
	}
}

func (c *ProductController) Create(ctx context.Context, req ProductCreateRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}
	return c.uc.Create(ctx, req.toEntity())
}

type ProductUpdateRequest struct {
	ID int64 `json:"id"`
	ProductAttr
}

func (r ProductUpdateRequest) toEntity() domain.Product {
	return domain.Product{
		ID:    r.ID,
		Name:  r.Name,
		Desc:  r.Desc,
		Price: r.Price,
	}
}

func (r ProductUpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required.Error("id is required")),
		validation.Field(&r.ProductAttr),
	)
}

func (c *ProductController) Update(ctx context.Context, req ProductUpdateRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	_, err := c.repo.FindOne(ctx, req.ID)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.WithError(err)
	} else if err != nil {
		return err
	}

	return c.uc.Update(ctx, req.toEntity())
}

func (c *ProductController) Delete(ctx context.Context, id int64) error {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	product, err := c.repo.FindOne(ctx, id)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.WithError(err)
	} else if err != nil {
		return err
	}

	return c.uc.Delete(ctx, *product)
}

func (c *ProductController) Detail(ctx context.Context, id int64) (*domain.Product, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	product, err := c.uc.Detail(ctx, id)
	if repository.IsNotFound(err) {
		return nil, berr.ErrResourceNotFound.WithError(err)
	} else if err != nil {
		return nil, err
	}

	return product, nil
}

type ProductListRequest struct {
	Keyword string
}

func (c *ProductController) List(ctx context.Context, req ProductListRequest) ([]*domain.Product, error) {
	param := usecase.ProductListParam{Keyword: req.Keyword}
	return c.uc.List(ctx, param)
}
