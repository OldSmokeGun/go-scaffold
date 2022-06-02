package greet

import (
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/service/greet"
	"go.uber.org/zap"
)

type HandlerInterface interface {
	Hello(ctx *gin.Context)
}

var _ HandlerInterface = (*Handler)(nil)

type Handler struct {
	logger  *zap.Logger
	service greet.ServiceInterface
}

func NewHandler(
	logger *zap.Logger,
	service greet.ServiceInterface,
) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}
