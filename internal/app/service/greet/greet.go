package greet

import (
	"go-scaffold/internal/app/global"
	"go.uber.org/zap"
)

type Service interface {
	Hello(param *HelloParam) (string, error)
}

type service struct {
	logger *zap.Logger
}

func New() *service {
	return &service{
		logger: global.Logger(),
	}
}
