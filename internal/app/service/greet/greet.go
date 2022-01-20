package greet

import (
	"go-scaffold/internal/app/global"
	"go.uber.org/zap"
)

type Interface interface {
	Hello(param *HelloParam) (string, error)
}

type service struct {
	Logger *zap.Logger
}

func New() *service {
	return &service{
		Logger: global.Logger(),
	}
}
