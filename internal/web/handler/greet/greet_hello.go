package greet

import (
	"gin-scaffold/internal/web/logic/greet"
	"gin-scaffold/internal/web/pkg/binx"
	"gin-scaffold/internal/web/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
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

func (w *handler) Hello(ctx *gin.Context) {
	req := new(HelloReq)
	if !binx.ShouldBindQuery(ctx, req) {
		return
	}

	param := &greet.HelloParam{
		Name: req.Name,
	}

	logic := greet.NewLogic(ctx)
	greetString, err := logic.Hello(param)
	if err != nil {
		ctx.JSON(http.StatusOK, response.ServerError.WithMsg(err.Error()))
	}

	resp := HelloResp{
		Msg: greetString,
	}

	ctx.JSON(http.StatusOK, response.Success.WithData(resp))
}
