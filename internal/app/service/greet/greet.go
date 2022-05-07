package greet

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type ServiceInterface interface {
	Hello(ctx context.Context, request HelloRequest) (*HelloResponse, error)
}

var _ ServiceInterface = (*Service)(nil)

type Service struct {
	logger *log.Helper
}

func NewService(logger log.Logger) *Service {
	return &Service{
		logger: log.NewHelper(logger),
	}
}
