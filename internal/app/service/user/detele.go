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
	user, err := s.Repository.FindOneByID(
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

	if err = s.Repository.Delete(user); err != nil {
		s.Logger.Error(err.Error())
		return ErrDataDeleteFailed
	}

	return nil
}
