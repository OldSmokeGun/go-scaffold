package greet

//go:generate mockgen -source=greet.go -destination=greet_mock.go -package=greet -mock_names=Interface=MockService

import (
	"context"
	"go-scaffold/internal/app/global"
	"go.uber.org/zap"
)

type Interface interface {
	Hello(ctx context.Context, param *HelloParam) (string, error)
}

type service struct {
	Logger *zap.Logger
}

func New() *service {
	return &service{
		Logger: global.Logger(),
	}
}
