package greet

import (
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/global"
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
