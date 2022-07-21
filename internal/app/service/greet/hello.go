package greet

import (
	"context"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	errorsx "go-scaffold/internal/app/pkg/errors"
)

type HelloRequest struct {
	Name string `json:"name" form:"name"`
}

func (r HelloRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required.Error("名称不能为空")),
	)
}

type HelloResponse struct {
	Msg string `json:"msg"`
}

func (s *Service) Hello(ctx context.Context, req HelloRequest) (*HelloResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errorsx.ValidateError().WithMessage(err.Error())
	}

	return &HelloResponse{
		Msg: fmt.Sprintf("Hello, %s", req.Name),
	}, nil
}
