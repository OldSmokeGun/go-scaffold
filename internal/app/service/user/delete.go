package user

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"go-scaffold/internal/app/model"
	errorsx "go-scaffold/internal/app/pkg/errors"
	"gorm.io/gorm"
)

// DeleteRequest 删除用户请求参数
type DeleteRequest struct {
	Id uint64 `json:"id" uri:"id"`
}

func (r DeleteRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Id, validation.Required.Error("id 不能为空")),
	)
}

// Delete 删除用户
func (s *Service) Delete(ctx context.Context, req DeleteRequest) error {
	if err := req.Validate(); err != nil {
		return errorsx.ValidateError().WithMessage(err.Error())
	}

	m, err := s.repo.FindOneById(
		ctx,
		req.Id,
		[]string{"*"},
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorsx.ResourceNotFound().WithMessage(model.ErrDataNotFound.Error())
		}
		s.logger.Errorf("%s: %s", model.ErrDataQueryFailed, err)
		return errorsx.ServerError().WithMessage(model.ErrDataQueryFailed.Error())
	}

	if err = s.repo.Delete(ctx, m); err != nil {
		s.logger.Errorf("%s: %s", model.ErrDataDeleteFailed, err)
		return errorsx.ServerError().WithMessage(model.ErrDataDeleteFailed.Error())
	}

	return nil
}
