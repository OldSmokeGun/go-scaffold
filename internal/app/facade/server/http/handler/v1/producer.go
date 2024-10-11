package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"go-scaffold/internal/app/controller"
	httperr "go-scaffold/internal/app/facade/server/http/pkg/errors"
)

type ProducerHandler struct {
	controller *controller.ProducerController
}

func NewProducerHandler(controller *controller.ProducerController) *ProducerHandler {
	return &ProducerHandler{controller}
}

type ProducerExampleRequest struct {
	Msg string `json:"msg"`
}

// Example 示例方法
//
//	@Router			/v1/producer/example [post]
//	@Summary		示例接口
//	@Description	示例接口
//	@Tags			示例
//	@Accept			json
//	@Produce		json
//	@Param			data	body		ProducerExampleRequest		true	"生产者消息"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *ProducerHandler) Example(ctx echo.Context) error {
	req := new(controller.ProducerExampleRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error")
	}

	if err := h.controller.Example(ctx.Request().Context(), *req); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}
