package greet

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"go-scaffold/internal/app/service/v1/greet"
	"go.uber.org/zap"
)

type Handler interface {
	Hello(ctx *gin.Context)
}

type handler struct {
	logger  *log.Helper
	zlogger *zap.Logger
	service *greet.Service
}

func New(logger log.Logger, zlogger *zap.Logger, service *greet.Service) Handler {
	return &handler{
		logger:  log.NewHelper(logger),
		zlogger: zlogger,
		service: service,
	}
}
