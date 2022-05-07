package greet

import (
	"github.com/gin-gonic/gin"
	errorsx "go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/service/greet"
	"go-scaffold/internal/app/transport/http/pkg/response"
)

type HelloRequest struct {
	greet.HelloRequest
}

// Hello 示例方法
// @Router       /v1/greet [get]
// @Summary      示例接口
// @Description  示例接口
// @Tags         示例
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        name  query     string                                     true  "名称"  format(string)  default(Tom)
// @Success      200   {object}  example.Success{data=greet.HelloResponse}  "成功响应"
// @Failure      500   {object}  example.ServerError                        "服务器出错"
// @Failure      400   {object}  example.ClientError                        "客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
// @Failure      401   {object}  example.Unauthorized                       "登陆失效"
// @Failure      403   {object}  example.PermissionDenied                   "没有权限"
// @Failure      404   {object}  example.ResourceNotFound                   "资源不存在"
// @Failure      429   {object}  example.TooManyRequest                     "请求过于频繁"
func (h *Handler) Hello(ctx *gin.Context) {
	req := new(HelloRequest)
	if err := ctx.ShouldBindQuery(req); err != nil {
		h.logger.Error(err)
		response.Error(ctx, errorsx.ValidateError())
		return
	}

	ret, err := h.service.Hello(ctx.Request.Context(), req.HelloRequest)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, response.WithData(ret))
	return
}
