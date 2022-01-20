package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go-scaffold/internal/app/rest/pkg/bindx"
	"go-scaffold/internal/app/rest/pkg/responsex"
	"go-scaffold/internal/app/service/user"
)

type (
	DetailReq struct {
		ID uint `uri:"id" binding:"required"` // 用户 ID
	}

	DetailResp struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`  // 名称
		Age   int8   `json:"age"`   // 年龄
		Phone string `json:"phone"` // 电话
	}
)

func (DetailReq) ErrorMessage() map[string]string {
	return map[string]string{
		"ID.required": "用户 ID 不能为空",
	}
}

// Detail 用户详情
// @Router       /v1/user/{id} [get]
// @Summary      用户详情
// @Description  用户详情
// @Tags         用户
// @Accept       plain
// @Produce      json
// @Param        id   path      integer                           true  "用户 ID"  format(uint)  minimum(1)
// @Success      200  {object}  example.Success{data=DetailResp}  "成功响应"
// @Failure      500  {object}  example.ServerError               "服务器出错"
// @Failure      400  {object}  example.ClientError               "客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
// @Failure      401  {object}  example.Unauthorized              "登陆失效"
// @Failure      403  {object}  example.PermissionDenied          "没有权限"
// @Failure      404  {object}  example.ResourceNotFound          "资源不存在"
// @Failure      429  {object}  example.TooManyRequest            "请求过于频繁"
func (h *handler) Detail(ctx *gin.Context) {
	req := new(DetailReq)
	if err := bindx.ShouldBindUri(ctx, req); err != nil {
		h.Logger.Error(err.Error())
		return
	}

	param := new(user.DetailParam)
	if err := copier.Copy(param, req); err != nil {
		h.Logger.Error(err.Error())
		responsex.ServerError(ctx)
		return
	}

	result, err := h.Service.Detail(param)
	if err != nil {
		responsex.ServerError(ctx, responsex.WithMsg(err.Error()))
		return
	}

	data := new(DetailResp)
	if err := copier.Copy(data, result); err != nil {
		h.Logger.Error(err.Error())
		responsex.ServerError(ctx)
		return
	}

	responsex.Success(ctx, responsex.WithData(data))
	return
}
