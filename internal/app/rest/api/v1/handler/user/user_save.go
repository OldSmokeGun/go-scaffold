package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go-scaffold/internal/app/rest/pkg/bindx"
	"go-scaffold/internal/app/rest/pkg/responsex"
	"go-scaffold/internal/app/service/user"
)

type SaveReq struct {
	ID    uint   `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`        // 名称
	Age   int8   `json:"age" binding:"min=1"`            // 年龄
	Phone string `json:"phone" binding:"required,phone"` // 电话
}

func (SaveReq) ErrorMessage() map[string]string {
	return map[string]string{
		"ID.required":    "用户 ID 不能为空",
		"Name.required":  "名称不能为空",
		"Age.min":        "年龄不能小于 {min}",
		"Phone.required": "手机号码不能为空",
		"Phone.phone":    "手机号码格式错误",
	}
}

// Save 更新用户
// @Router       /v1/user [put]
// @Summary      更新用户
// @Description  更新用户
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        user_info  body      SaveReq                   true  "用户信息"  format(string)
// @Success      200        {object}  example.Success           "成功响应"
// @Failure      500        {object}  example.ServerError       "服务器出错"
// @Failure      400        {object}  example.ClientError       "客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
// @Failure      401        {object}  example.Unauthorized      "登陆失效"
// @Failure      403        {object}  example.PermissionDenied  "没有权限"
// @Failure      404        {object}  example.ResourceNotFound  "资源不存在"
// @Failure      429        {object}  example.TooManyRequest    "请求过于频繁"
func (h *handler) Save(ctx *gin.Context) {
	req := new(SaveReq)
	if err := bindx.ShouldBindJSON(ctx, req); err != nil {
		h.logger.Error(err.Error())
		return
	}

	param := new(user.SaveParam)
	if err := copier.Copy(param, req); err != nil {
		h.logger.Error(err.Error())
		responsex.ServerError(ctx)
		return
	}

	err := h.service.Save(param)
	if err != nil {
		responsex.ServerError(ctx, responsex.WithMsg(err.Error()))
		return
	}

	responsex.Success(ctx)
	return
}
