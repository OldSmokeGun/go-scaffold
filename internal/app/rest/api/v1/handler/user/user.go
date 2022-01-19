package user

import (
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/service/user"
	"go.uber.org/zap"
)

type Handler interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Detail(ctx *gin.Context)
	List(ctx *gin.Context)
	Save(ctx *gin.Context)
}

type handler struct {
	logger  *zap.Logger
	service user.Service
}

func New() *handler {
	return &handler{
		logger:  global.Logger(),
		service: user.New(),
	}
}
