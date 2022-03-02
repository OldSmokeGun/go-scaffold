package greet

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
)

type Handler interface {
	Default(cmd *cobra.Command, args []string)
	To(cmd *cobra.Command, args []string)
}

type handler struct {
	logger *log.Helper
}

func NewHandler(logger log.Logger) Handler {
	return &handler{
		logger: log.NewHelper(logger),
	}
}
