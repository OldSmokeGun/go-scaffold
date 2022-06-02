package user

import (
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/service/user"
	"go.uber.org/zap"
)

type HandlerInterface interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Detail(ctx *gin.Context)
	List(ctx *gin.Context)
}

var _ HandlerInterface = (*Handler)(nil)

type Handler struct {
	logger  *zap.Logger
	service user.ServiceInterface
}

func NewHandler(
	logger *zap.Logger,
	service user.ServiceInterface,
) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}
