package user

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	pb "go-scaffold/internal/app/api/v1/user"
)

type HandlerInterface interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Detail(ctx *gin.Context)
	List(ctx *gin.Context)
}

type Handler struct {
	logger  *log.Helper
	service pb.UserServer
}

func NewHandler(logger log.Logger, service pb.UserServer) *Handler {
	return &Handler{
		logger:  log.NewHelper(logger),
		service: service,
	}
}
