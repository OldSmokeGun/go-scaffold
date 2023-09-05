package handler

import "log/slog"

type ExampleMessage struct {
	Msg string `json:"msg"`
}

type ExampleHandler struct {
	logger *slog.Logger
}

func NewExampleHandler(logger *slog.Logger) *ExampleHandler {
	return &ExampleHandler{logger}
}

func (h *ExampleHandler) Handle(message ExampleMessage) error {
	h.logger.Info("receive example message: " + message.Msg)
	return nil
}
