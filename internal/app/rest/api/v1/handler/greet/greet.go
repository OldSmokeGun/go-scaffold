package greet

import (
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/service/greet"
	"go.uber.org/zap"
)

type Handler interface {
	Hello(ctx *gin.Context)
}

type handler struct {
	logger  *zap.Logger
	service greet.Service
}

func New() *handler {
	return &handler{
		logger:  global.Logger(),
		service: greet.New(),
	}
}
