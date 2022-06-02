package greet

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Handler interface {
	Default(cmd *cobra.Command, args []string)
	To(cmd *cobra.Command, args []string)
}

type handler struct {
	logger *zap.Logger
}

func NewHandler(logger *zap.Logger) Handler {
	return &handler{
		logger: logger,
	}
}
