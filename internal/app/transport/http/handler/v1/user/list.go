package user

import (
	"github.com/gin-gonic/gin"
	pb "go-scaffold/internal/app/api/v1/user"
	"go-scaffold/internal/app/pkg/responsex"
	"go-scaffold/internal/app/transport/http/pkg/bindx"
)

type ListRequest struct {
	pb.ListRequest
}

func (*ListRequest) Message() map[string]string {
	return nil
}

// List 用户列表
// @Router       /v1/users [get]
// @Summary      用户列表
// @Description  用户列表
// @Tags         用户
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        keyword  query     string                                 false  "查询字符串"  format(string)
// @Success      200      {object}  example.Success{data=pb.ListResponse}  "成功响应"
// @Failure      500      {object}  example.ServerError                    "服务器出错"
// @Failure      400      {object}  example.ClientError                    "客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
// @Failure      401      {object}  example.Unauthorized                   "登陆失效"
// @Failure      403      {object}  example.PermissionDenied               "没有权限"
// @Failure      404      {object}  example.ResourceNotFound               "资源不存在"
// @Failure      429      {object}  example.TooManyRequest                 "请求过于频繁"
func (h *Handler) List(ctx *gin.Context) {
	req := new(ListRequest)
	if err := bindx.ShouldBindQuery(ctx, req); err != nil {
		h.logger.Error(err.Error())
		return
	}

	ret, err := h.service.List(ctx.Request.Context(), &req.ListRequest)
	if err != nil {
		responsex.ServerError(ctx, responsex.WithMsg(err.Error()))
		return
	}

	responsex.Success(ctx, responsex.WithData(ret))
	return
}
