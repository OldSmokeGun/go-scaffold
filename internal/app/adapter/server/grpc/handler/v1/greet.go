package v1

import (
	"context"

	v1 "go-scaffold/internal/app/adapter/server/grpc/api/v1"
	"go-scaffold/internal/app/adapter/server/grpc/pkg/errors"
	"go-scaffold/internal/app/controller"

	"golang.org/x/exp/slog"
)

// GreetHandler 示例处理器
type GreetHandler struct {
	v1.UnimplementedGreetServer
	logger     *slog.Logger
	controller *controller.GreetController
}

// NewGreetHandler 构造示例处理器
func NewGreetHandler(
	logger *slog.Logger,
	controller *controller.GreetController,
) *GreetHandler {
	return &GreetHandler{
		logger:     logger,
		controller: controller,
	}
}

// Hello 示例方法
func (h *GreetHandler) Hello(ctx context.Context, req *v1.GreetHelloRequest) (*v1.GreetHelloResponse, error) {
	controllerReq := controller.HelloRequest{Name: req.Name}

	ret, err := h.controller.Hello(ctx, controllerReq)
	if err != nil {
		h.logger.Error("call controller.GreetController.Hello method error", err)
		return nil, errors.Wrap(err)
	}

	resp := &v1.GreetHelloResponse{
		Msg: ret.Msg,
	}

	return resp, nil
}
