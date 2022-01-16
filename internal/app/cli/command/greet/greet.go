package greet

import (
	"github.com/spf13/cobra"
	"go-scaffold/internal/app/global"
	"go.uber.org/zap"
)

type Handler interface {
	Default(cmd *cobra.Command, args []string)
}

type handler struct {
	logger *zap.Logger
}

func NewHandler() *handler {
	return &handler{
		logger: global.Logger(),
	}
}
