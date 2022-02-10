package user

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"go-scaffold/internal/app/rest/pkg/responsex"
	"gorm.io/gorm"
)

type (
	DetailParam struct {
		ID uint // 用户 ID
	}

	DetailResult struct {
		ID    uint
		Name  string
		Age   int8
		Phone string
	}
)

// Detail 用户详情
func (s *service) Detail(ctx context.Context, param *DetailParam) (*DetailResult, error) {
	user, err := s.Repository.FindOneByID(
		context.TODO(),
		param.ID,
		[]string{"*"},
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotExist
		}
		s.Logger.Error(err.Error())
		return nil, ErrDataQueryFailed
	}

	result := new(DetailResult)
	if err = copier.Copy(result, user); err != nil {
		s.Logger.Error(err.Error())
		return nil, errors.New(responsex.ServerErrorCode.String())
	}

	return result, nil
}
