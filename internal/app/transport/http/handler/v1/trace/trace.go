package trace

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"go-scaffold/internal/app/component/trace"
	"go-scaffold/internal/app/config"
)

type HandlerInterface interface {
	Example(ctx *gin.Context)
}

type Handler struct {
	logger *log.Helper
	cm     *config.Config
	trace  *trace.Tracer
}

func NewHandler(logger log.Logger, cm *config.Config, trace *trace.Tracer) *Handler {
	return &Handler{
		logger: log.NewHelper(logger),
		cm:     cm,
		trace:  trace,
	}
}
