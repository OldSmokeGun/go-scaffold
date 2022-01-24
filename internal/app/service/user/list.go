package user

import (
	"errors"
	"github.com/jinzhu/copier"
	"go-scaffold/internal/app/rest/pkg/responsex"
)

type ListParam struct {
	Keyword string // 查询字符串
}

type ListResult []*struct {
	ID    uint
	Name  string
	Age   int8
	Phone string
}

// List 用户列表
func (s *service) List(param *ListParam) (ListResult, error) {
	users, err := s.Repository.FindByKeyword(
		[]string{"*"},
		param.Keyword,
		"updated_at DESC",
	)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrDataQueryFailed
	}

	var result ListResult
	if err = copier.Copy(&result, users); err != nil {
		s.Logger.Error(err.Error())
		return nil, errors.New(responsex.ServerErrorCode.String())
	}

	return result, nil
}
