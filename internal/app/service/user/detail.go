package user

import (
	"errors"
	"go-scaffold/internal/app/model"
	"gorm.io/gorm"
)

type DetailParam struct {
	ID uint // 用户 ID
}

// Detail 用户详情
func (s *service) Detail(param *DetailParam) (*model.User, error) {
	user, err := s.Repository.FindOneByID(
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

	return user, nil
}
