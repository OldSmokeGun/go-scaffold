package user

import (
	"errors"
	"gorm.io/gorm"
)

type DeleteParam struct {
	ID uint // 用户 ID
}

// Delete 删除用户
func (s *service) Delete(param *DeleteParam) error {
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

	if err = s.repository.Delete(user); err != nil {
		s.logger.Error(err.Error())
		return errors.New("数据删除失败")
	}

	return nil
}
