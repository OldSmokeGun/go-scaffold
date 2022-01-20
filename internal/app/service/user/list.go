package user

import (
	"go-scaffold/internal/app/model"
)

type ListParam struct {
	Keyword string // 查询字符串
}

// List 用户列表
func (s *service) List(param *ListParam) ([]*model.User, error) {
	users, err := s.Repository.FindByKeyword(
		[]string{"*"},
		param.Keyword,
		"updated_at DESC",
	)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrDataQueryFailed
	}

	return users, nil
}
