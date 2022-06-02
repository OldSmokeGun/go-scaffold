package greet

import (
	"context"
	"go.uber.org/zap"
)

type ServiceInterface interface {
	Hello(ctx context.Context, request HelloRequest) (*HelloResponse, error)
}

var _ ServiceInterface = (*Service)(nil)

type Service struct {
	logger *zap.Logger
}

func NewService(logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
	}
}
