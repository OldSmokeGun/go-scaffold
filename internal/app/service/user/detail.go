package user

import (
	"context"
	"errors"
	"go-scaffold/internal/app/model"
	errorsx "go-scaffold/internal/app/pkg/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

// DetailRequest 用户详情请求参数
type DetailRequest struct {
	Id uint64 `json:"id" uri:"id"`
}

func (c DetailRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Id, validation.Required.Error("id 不能为空")),
	)
}

// DetailResponse 用户详情响应数据
type DetailResponse struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Age   int8   `json:"age"`
	Phone string `json:"phone"`
}

// Detail 用户详情
func (s *Service) Detail(ctx context.Context, req DetailRequest) (*DetailResponse, error) {
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

	resp := &DetailResponse{
		Id:    m.Id,
		Name:  m.Name,
		Age:   m.Age,
		Phone: m.Phone,
	}

	return resp, nil
}
