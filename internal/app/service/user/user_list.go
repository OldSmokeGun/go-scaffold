package user

import (
	"errors"
	"go-scaffold/internal/app/model"
)

type ListParam struct {
	Keyword string // 查询字符串
}

// List 用户列表
func (s *service) List(param *ListParam) ([]*model.User, error) {
	users, err := s.repository.FindByKeyword(
		[]string{"*"},
		param.Keyword,
		"updated_at DESC",
	)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, errors.New("数据查询失败")
	}

	return users, nil
}
