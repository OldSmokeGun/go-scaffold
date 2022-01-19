package user

import (
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/model"
	"go-scaffold/internal/app/repository/user"
	"go.uber.org/zap"
)

type Service interface {
	Create(param *CreateParam) error
	List(param *ListParam) ([]*model.User, error)
	Detail(param *DetailParam) (*model.User, error)
	Save(param *SaveParam) error
	Delete(param *DeleteParam) error
}

type service struct {
	logger     *zap.Logger
	repository user.Repository
}

func New() *service {
	return &service{
		logger:     global.Logger(),
		repository: user.New(),
	}
}
