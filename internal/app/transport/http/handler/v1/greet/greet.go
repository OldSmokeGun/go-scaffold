package greet

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/service/v1/greet"
	"go.uber.org/zap"
)

type Handler interface {
	Hello(ctx *gin.Context)
}

type handler struct {
	logger  *log.Helper
	zLogger *zap.Logger
	cm      *config.Config
	service *greet.Service
}

func New(logger log.Logger, zLogger *zap.Logger, cm *config.Config, service *greet.Service) Handler {
	return &handler{
		logger:  log.NewHelper(logger),
		zLogger: zLogger,
		cm:      cm,
		service: service,
	}
}
