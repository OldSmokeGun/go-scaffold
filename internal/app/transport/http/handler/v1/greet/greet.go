package greet

import (
	"go-scaffold/internal/app/service/greet"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type HandlerInterface interface {
	Hello(ctx *gin.Context)
}

var _ HandlerInterface = (*Handler)(nil)

type Handler struct {
	logger  *log.Helper
	service greet.ServiceInterface
}

func NewHandler(
	logger log.Logger,
	service greet.ServiceInterface,
) *Handler {
	return &Handler{
		logger:  log.NewHelper(logger),
		service: service,
	}
}
