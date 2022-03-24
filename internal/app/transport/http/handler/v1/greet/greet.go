package greet

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	pb "go-scaffold/internal/app/api/scaffold/v1/greet"
	"go-scaffold/internal/app/config"
	"go.uber.org/zap"
)

type HandlerInterface interface {
	Hello(ctx *gin.Context)
}

type Handler struct {
	logger  *log.Helper
	zLogger *zap.Logger
	conf    *config.Config
	service pb.GreetServer
}

func NewHandler(logger log.Logger, zLogger *zap.Logger, conf *config.Config, service pb.GreetServer) *Handler {
	return &Handler{
		logger:  log.NewHelper(logger),
		zLogger: zLogger,
		conf:    conf,
		service: service,
	}
}
