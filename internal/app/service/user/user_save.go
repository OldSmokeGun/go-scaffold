package user

import (
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
func (s *service) Save(param *SaveParam) error {
	user, err := s.repository.FindOneByID(
		param.ID,
		[]string{"*"},
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		s.logger.Error(err.Error())
		return errors.New("数据查询失败")
	}

	if err = copier.Copy(user, param); err != nil {
		s.logger.Error(err.Error())
		return errors.New(responsex.ServerErrorCode.String())
	}

	if _, err = s.repository.Save(user); err != nil {
		s.logger.Error(err.Error())
		return errors.New("数据保存失败")
	}

	return nil
}
