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
	user, err := s.repository.FindOneByID(
		param.ID,
		[]string{"*"},
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		s.logger.Error(err.Error())
		return nil, errors.New("数据查询失败")
	}

	return user, nil
}
