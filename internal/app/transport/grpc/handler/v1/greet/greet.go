package greet

import (
	"github.com/go-kratos/kratos/v2/log"
	pb "go-scaffold/internal/app/api/scaffold/v1/greet"
	"go-scaffold/internal/app/service/greet"
)

var _ pb.GreetServer = (*Handler)(nil)

type Handler struct {
	pb.UnimplementedGreetServer
	logger  *log.Helper
	service *greet.Service
}

func NewHandler(
	logger log.Logger,
	service *greet.Service,
) *Handler {
	return &Handler{
		logger:  log.NewHelper(logger),
		service: service,
	}
}
