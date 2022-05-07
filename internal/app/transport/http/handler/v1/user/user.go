package user

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"go-scaffold/internal/app/service/user"
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
	logger  *log.Helper
	service user.ServiceInterface
}

func NewHandler(
	logger log.Logger,
	service user.ServiceInterface,
) *Handler {
	return &Handler{
		logger:  log.NewHelper(logger),
		service: service,
	}
}
