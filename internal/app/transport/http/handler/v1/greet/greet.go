package greet

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/service/v1/greet"
	"go.uber.org/zap"
)

type HandlerInterface interface {
	Hello(ctx *gin.Context)
}

type Handler struct {
	logger  *log.Helper
	zLogger *zap.Logger
	cm      *config.Config
	service *greet.Service
}

func NewHandler(logger log.Logger, zLogger *zap.Logger, cm *config.Config, service *greet.Service) *Handler {
	return &Handler{
		logger:  log.NewHelper(logger),
		zLogger: zLogger,
		cm:      cm,
		service: service,
	}
}
