package user

import (
	"errors"
	"github.com/jinzhu/copier"
	"go-scaffold/internal/app/model"
	"go-scaffold/internal/app/rest/pkg/responsex"
)

type CreateParam struct {
	Name  string // 名称
	Age   int8   // 年龄
	Phone string // 电话
}

// Create 创建用户
func (s *service) Create(param *CreateParam) error {
	m := new(model.User)
	if err := copier.Copy(m, param); err != nil {
		s.logger.Error(err.Error())
		return errors.New(responsex.ServerErrorCode.String())
	}

	if _, err := s.repository.Create(m); err != nil {
		s.logger.Error(err.Error())
		return errors.New("数据保存失败")
	}

	return nil
}
