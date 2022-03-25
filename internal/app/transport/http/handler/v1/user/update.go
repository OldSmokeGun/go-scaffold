package user

import (
	"github.com/gin-gonic/gin"
	pb "go-scaffold/internal/app/api/scaffold/v1/user"
	"go-scaffold/internal/app/transport/http/pkg/bindx"
	"go-scaffold/internal/app/transport/http/pkg/responsex"
	"strconv"
)

type UpdateRequest struct {
	pb.UpdateRequest
}

func (*UpdateRequest) Message() map[string]string {
	return map[string]string{
		"UpdateRequest.Id.required":    "用户 id 不能为空",
		"UpdateRequest.Name.required":  "名称不能为空",
		"UpdateRequest.Age.min":        "年龄必须大于 {min}",
		"UpdateRequest.Phone.required": "手机号码不能为空",
		"UpdateRequest.Phone.phone":    "手机号码格式错误",
	}
}

// Update 更新用户
// @Router       /v1/user/{id} [put]
// @Summary      更新用户
// @Description  更新用户
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        id         path      integer                   true  "用户 id"  format(uint)  minimum(1)
// @Param        user_info  body      CreateRequest             true  "用户信息"   format(string)
// @Success      200        {object}  example.Success           "成功响应"
// @Failure      500        {object}  example.ServerError       "服务器出错"
// @Failure      400        {object}  example.ClientError       "客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
// @Failure      401        {object}  example.Unauthorized      "登陆失效"
// @Failure      403        {object}  example.PermissionDenied  "没有权限"
// @Failure      404        {object}  example.ResourceNotFound  "资源不存在"
// @Failure      429        {object}  example.TooManyRequest    "请求过于频繁"
func (h *Handler) Update(ctx *gin.Context) {
	var err error

	req := new(UpdateRequest)
	req.Id, err = strconv.ParseUint(ctx.Params.ByName("id"), 10, 64)
	if err != nil {
		responsex.ValidateError(ctx, responsex.WithMsg("参数解析失败"))
		h.logger.Error(err)
		return
	}

	if err = bindx.ShouldBindJSON(ctx, req); err != nil {
		h.logger.Error(err)
		return
	}

	_, err = h.service.Update(ctx.Request.Context(), &req.UpdateRequest)
	if err != nil {
		responsex.ServerError(ctx, responsex.WithMsg(err.Error()))
		return
	}

	responsex.Success(ctx)
	return
}
