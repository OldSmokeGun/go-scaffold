package greet

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go-scaffold/internal/app/logic/greet"
	"go-scaffold/internal/app/rest/pkg/bindx"
	"go-scaffold/internal/app/rest/pkg/responsex"
)

type (
	HelloReq struct {
		Name string `form:"name" binding:"required"`
	}

	HelloResp struct {
		Msg string `json:"msg"`
	}
)

func (HelloReq) ErrorMessage() map[string]string {
	return map[string]string{
		"Name.required": "名称不能为空",
	}
}

// Hello 示例方法
// @Router       /greet [get]
// @Summary      示例接口
// @Description  示例接口
// @Tags         实例
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        name  query     string                       true  "名称"  format(string)  default(Tom)
// @Success      200   {object}  api.Success{data=HelloResp}  "成功响应"
// @Failure      500   {object}  api.ServerError              "服务器出错"
// @Failure      400   {object}  api.ClientError              "客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
// @Failure      401   {object}  api.Unauthorized             "登陆失效"
// @Failure      403   {object}  api.PermissionDenied         "没有权限"
// @Failure      404   {object}  api.ResourceNotFound         "资源不存在"
// @Failure      429   {object}  api.TooManyRequest           "请求过于频繁"
func (h *handler) Hello(ctx *gin.Context) {
	req := new(HelloReq)
	if err := bindx.ShouldBindQuery(ctx, req); err != nil {
		h.logger.Error(err.Error())
		return
	}

	param := new(greet.HelloParam)
	if err := copier.Copy(param, req); err != nil {
		h.logger.Error(err.Error())
		responsex.ServerError(ctx)
		return
	}

	logic := greet.NewLogic(ctx)
	greetString, err := logic.Hello(param)
	if err != nil {
		responsex.ServerError(ctx, responsex.WithMsg(err.Error()))
		return
	}

	resp := HelloResp{
		Msg: greetString,
	}

	responsex.Success(ctx, responsex.WithData(resp))
	return
}
