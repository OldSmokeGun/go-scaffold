package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go-scaffold/internal/app/api/v1/user"
	"go-scaffold/internal/app/transport/http/pkg/responsex"
)

type (
	ListReq struct {
		Keyword string `form:"keyword"` // 查询字符串
	}

	ListItem struct {
		Id    uint64 `json:"id"`
		Name  string `json:"name"`  // 名称
		Age   int32  `json:"age"`   // 年龄
		Phone string `json:"phone"` // 电话
	}

	ListResp []*ListItem
)

func (ListReq) ErrorMessage() map[string]string {
	return nil
}

// List 用户列表
// @Router       /v1/users [get]
// @Summary      用户列表
// @Description  用户列表
// @Tags         用户
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        keyword  query     string                          false  "查询字符串"  format(string)
// @Success      200      {object}  example.Success{data=ListResp}  "成功响应"
// @Failure      500      {object}  example.ServerError             "服务器出错"
// @Failure      400      {object}  example.ClientError             "客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
// @Failure      401      {object}  example.Unauthorized            "登陆失效"
// @Failure      403      {object}  example.PermissionDenied        "没有权限"
// @Failure      404      {object}  example.ResourceNotFound        "资源不存在"
// @Failure      429      {object}  example.TooManyRequest          "请求过于频繁"
func (h *handler) List(ctx *gin.Context) {
	req := &ListReq{
		Keyword: ctx.Query("keyword"),
	}

	param := new(user.ListRequest)
	if err := copier.Copy(param, req); err != nil {
		h.logger.Error(err.Error())
		responsex.ServerError(ctx)
		return
	}

	ret, err := h.service.List(ctx.Request.Context(), param)
	if err != nil {
		responsex.ServerError(ctx, responsex.WithMsg(err.Error()))
		return
	}

	var data ListResp
	for _, u := range ret.Users {
		data = append(data, &ListItem{
			Id:    u.Id,
			Name:  u.Name,
			Age:   u.Age,
			Phone: u.Phone,
		})
	}

	responsex.Success(ctx, responsex.WithData(data))
	return
}
