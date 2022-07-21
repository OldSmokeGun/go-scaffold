package user

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go-scaffold/internal/app/model"
	errorsx "go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/pkg/validator"
	"gorm.io/gorm"
)

// UpdateRequest 更新用户请求参数
type UpdateRequest struct {
	Id    uint64 `json:"id" uri:"id"`
	Name  string `json:"name"`
	Age   int8   `json:"age"`
	Phone string `json:"phone"`
}

func (r UpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Id, validation.Required.Error("id 不能为空")),
		validation.Field(&r.Name, validation.Required.Error("名称不能为空")),
		validation.Field(&r.Phone, validation.By(validator.IsMobilePhone)),
	)
}

// UpdateResponse 更新用户响应数据
type UpdateResponse struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Age   int8   `json:"age"`
	Phone string `json:"phone"`
}

// Update 更新用户
func (s *Service) Update(ctx context.Context, req UpdateRequest) (*UpdateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errorsx.ValidateError().WithMessage(err.Error())
	}

	m, err := s.repo.FindOneById(
		ctx,
		req.Id,
		[]string{"*"},
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorsx.ResourceNotFound().WithMessage(model.ErrDataNotFound.Error())
		}
		s.logger.Errorf("%s: %s", model.ErrDataQueryFailed, err)
		return nil, errorsx.ServerError().WithMessage(model.ErrDataQueryFailed.Error())
	}

	m = &model.User{
		BaseModel: m.BaseModel,
		Name:      req.Name,
		Age:       req.Age,
		Phone:     req.Phone,
	}

	if _, err = s.repo.Save(ctx, m); err != nil {
		s.logger.Errorf("%s: %s", model.ErrDataStoreFailed, err)
		return nil, errorsx.ServerError().WithMessage(model.ErrDataStoreFailed.Error())
	}

	resp := &UpdateResponse{
		Id:    m.Id,
		Name:  m.Name,
		Age:   m.Age,
		Phone: m.Phone,
	}

	return resp, nil
}
