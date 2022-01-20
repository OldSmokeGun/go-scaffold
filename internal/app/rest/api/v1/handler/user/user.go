package user

import (
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/service/user"
	"go.uber.org/zap"
)

type Interface interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Detail(ctx *gin.Context)
	List(ctx *gin.Context)
	Save(ctx *gin.Context)
}

type handler struct {
	Logger  *zap.Logger
	Service user.Interface
}

func New() *handler {
	return &handler{
		Logger:  global.Logger(),
		Service: user.New(),
	}
}
