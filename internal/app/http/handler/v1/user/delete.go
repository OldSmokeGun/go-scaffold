package user

import (
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/http/pkg/response"
	errorsx "go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/service/user"
)

// Delete 删除用户
// @Router       /v1/user/{id} [delete]
// @Summary      删除用户
// @Description  删除用户
// @Tags         用户
// @Accept       plain
// @Produce      json
// @Param        id   path      integer                   true  "用户 id"  format(uint)  minimum(1)
// @Success      200  {object}  example.Success           "成功响应"
// @Failure      500  {object}  example.ServerError       "服务器出错"
// @Failure      400  {object}  example.ClientError       "客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
// @Failure      401  {object}  example.Unauthorized      "登陆失效"
// @Failure      403  {object}  example.PermissionDenied  "没有权限"
// @Failure      404  {object}  example.ResourceNotFound  "资源不存在"
// @Failure      429  {object}  example.TooManyRequest    "请求过于频繁"
// @Security     Authorization
func (h *Handler) Delete(ctx *gin.Context) {
	req := new(user.DeleteRequest)
	if err := ctx.ShouldBindUri(req); err != nil {
		h.logger.Sugar().Error(err)
		response.Error(ctx, errorsx.ValidateError())
		return
	}

	err := h.service.Delete(ctx.Request.Context(), *req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
	return
}
