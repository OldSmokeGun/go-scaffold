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

// HelloRequest 请求参数
type HelloRequest struct {
	Name string `json:"name" query:"name"`
}

// Validate 验证参数
func (r HelloRequest) Validate() error {
	return errors.WithStack(validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required.Error("名称不能为空")),
	))
}

// HelloResponse 请求响应
type HelloResponse struct {
	Msg string
}

// Hello 示例方法
func (s *GreetController) Hello(ctx context.Context, req HelloRequest) (*HelloResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.WithStack(berr.ErrValidateError.WithMsg(err.Error()))
	}

	return &HelloResponse{
		Msg: fmt.Sprintf("Hello, %s", req.Name),
	}, nil
}
