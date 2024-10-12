package v1

import (
	"context"
	"log/slog"

	"go-scaffold/internal/app/controller"
	v1 "go-scaffold/internal/app/facade/server/grpc/api/v1"
	"go-scaffold/internal/app/facade/server/grpc/pkg/errors"
)

type GreetHandler struct {
	v1.UnimplementedGreetServer
	logger     *slog.Logger
	controller *controller.GreetController
}

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
	controllerReq := controller.GreetHelloRequest{Name: req.Name}

	ret, err := h.controller.Hello(ctx, controllerReq)
	if err != nil {
		h.logger.Error("call GreetController.Hello method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.GreetHelloResponse{
		Msg: ret.Msg,
	}, nil
}
