package user

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go-scaffold/internal/app/model"
	errorsx "go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/pkg/validator"
)

// CreateRequest 创建用户请求参数
type CreateRequest struct {
	Name  string `json:"name"`
	Age   int8   `json:"age"`
	Phone string `json:"phone"`
}

func (r CreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required.Error("名称不能为空")),
		validation.Field(&r.Phone, validation.By(validator.IsMobilePhone)),
	)
}

// CreateResponse 创建用户响应数据
type CreateResponse struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Age   int8   `json:"age"`
	Phone string `json:"phone"`
}

// Create 创建用户
func (s *Service) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errorsx.ValidateError().WithMessage(err.Error())
	}

	m := &model.User{
		Name:  req.Name,
		Age:   req.Age,
		Phone: req.Phone,
	}

	if _, err := s.repo.Create(ctx, m); err != nil {
		s.logger.Errorf("%s: %s", model.ErrDataStoreFailed, err)
		return nil, errorsx.ServerError().WithMessage(model.ErrDataStoreFailed.Error())
	}

	resp := &CreateResponse{
		Id:    m.Id,
		Name:  m.Name,
		Age:   m.Age,
		Phone: m.Phone,
	}

	return resp, nil
}
