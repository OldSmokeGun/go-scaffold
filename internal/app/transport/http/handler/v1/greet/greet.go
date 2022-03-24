package greet

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	pb "go-scaffold/internal/app/api/scaffold/v1/greet"
)

type HandlerInterface interface {
	Hello(ctx *gin.Context)
}

type Handler struct {
	logger  *log.Helper
	service pb.GreetServer
}

func NewHandler(
	logger log.Logger,
	service pb.GreetServer,
) *Handler {
	return &Handler{
		logger:  log.NewHelper(logger),
		service: service,
	}
}
