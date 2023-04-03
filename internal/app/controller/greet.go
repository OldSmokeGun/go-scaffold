package controller

import (
	"context"
	"fmt"

	berr "go-scaffold/internal/app/pkg/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

// GreetController 示例控制器
type GreetController struct{}

// NewGreetController 构造示例控制器
func NewGreetController() *GreetController {
	return &GreetController{}
}

// GreetHelloRequest 请求参数
type GreetHelloRequest struct {
	Name string `json:"name" query:"name"`
}

// Validate 验证参数
func (r GreetHelloRequest) Validate() error {
	return errors.WithStack(validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required.Error("名称不能为空")),
	))
}

// GreetHelloResponse 请求响应
type GreetHelloResponse struct {
	Msg string
}

// Hello 示例方法
func (c *GreetController) Hello(ctx context.Context, req GreetHelloRequest) (*GreetHelloResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.WithStack(berr.ErrValidateError.WithMsg(err.Error()))
	}

	return &GreetHelloResponse{
		Msg: fmt.Sprintf("Hello, %s", req.Name),
	}, nil
}
