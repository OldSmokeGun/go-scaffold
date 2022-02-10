package trace

import (
	"github.com/gin-gonic/gin"
	appconfig "go-scaffold/internal/app/config"
	"go-scaffold/internal/app/global"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Interface interface {
	Example(ctx *gin.Context)
}

type handler struct {
	Config *appconfig.Config
	Tracer oteltrace.Tracer
}

func New() *handler {
	return &handler{
		Config: global.Config(),
		Tracer: global.Tracer(),
	}
}
