package user

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"go-scaffold/internal/app/rest/pkg/responsex"
	"gorm.io/gorm"
)

type SaveParam struct {
	ID    uint
	Name  string // 名称
	Age   int8   // 年龄
	Phone string // 电话
}

// Save 更新用户
func (s *service) Save(ctx context.Context, param *SaveParam) error {
	user, err := s.Repository.FindOneByID(
		context.TODO(),
		param.ID,
		[]string{"*"},
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotExist
		}
		s.Logger.Error(err.Error())
		return ErrDataQueryFailed
	}

	if err = copier.Copy(user, param); err != nil {
		s.Logger.Error(err.Error())
		return errors.New(responsex.ServerErrorCode.String())
	}

	if _, err = s.Repository.Save(context.TODO(), user); err != nil {
		s.Logger.Error(err.Error())
		return ErrDataStoreFailed
	}

	return nil
}
