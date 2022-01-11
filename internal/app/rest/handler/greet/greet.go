package greet

import (
	"gin-scaffold/internal/app/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler interface {
	Hello(ctx *gin.Context)
}

type handler struct {
	logger *zap.Logger
}

func NewHandler() *handler {
	return &handler{
		logger: global.Logger(),
	}
}
