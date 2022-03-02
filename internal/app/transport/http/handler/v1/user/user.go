package user

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"go-scaffold/internal/app/service/v1/user"
)

type Handler interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Detail(ctx *gin.Context)
	List(ctx *gin.Context)
	Save(ctx *gin.Context)
}

type handler struct {
	logger  *log.Helper
	service *user.Service
}

func New(logger log.Logger, service *user.Service) Handler {
	return &handler{
		logger:  log.NewHelper(logger),
		service: service,
	}
}
