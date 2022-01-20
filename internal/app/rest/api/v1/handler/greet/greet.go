package greet

import (
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/service/greet"
	"go.uber.org/zap"
)

type Interface interface {
	Hello(ctx *gin.Context)
}

type handler struct {
	Logger  *zap.Logger
	Service greet.Interface
}

func New() *handler {
	return &handler{
		Logger:  global.Logger(),
		Service: greet.New(),
	}
}
