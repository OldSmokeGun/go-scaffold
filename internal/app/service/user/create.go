package user

import (
	"context"
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
func (s *service) Create(ctx context.Context, param *CreateParam) error {
	m := new(model.User)
	if err := copier.Copy(m, param); err != nil {
		s.Logger.Error(err.Error())
		return errors.New(responsex.ServerErrorCode.String())
	}

	if _, err := s.Repository.Create(context.TODO(), m); err != nil {
		s.Logger.Error(err.Error())
		return ErrDataStoreFailed
	}

	return nil
}
