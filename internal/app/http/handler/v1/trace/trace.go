package trace

import (
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/component/trace"
	"go-scaffold/internal/app/config"
	"go.uber.org/zap"
)

type HandlerInterface interface {
	Example(ctx *gin.Context)
}

var _ HandlerInterface = (*Handler)(nil)

type Handler struct {
	logger *zap.Logger
	conf   *config.Config
	trace  *trace.Tracer
}

func NewHandler(
	logger *zap.Logger,
	conf *config.Config,
	trace *trace.Tracer,
) *Handler {
	return &Handler{
		logger: logger,
		conf:   conf,
		trace:  trace,
	}
}
