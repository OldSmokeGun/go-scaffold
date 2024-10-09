package controller

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"

	berr "go-scaffold/internal/errors"
)

type GreetController struct{}

func NewGreetController() *GreetController {
	return &GreetController{}
}

type GreetHelloRequest struct {
	Name string `json:"name" query:"name"`
}

func (r GreetHelloRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required.Error("name is required")),
	)
}

type GreetHelloResponse struct {
	Msg string
}

func (c *GreetController) Hello(ctx context.Context, req GreetHelloRequest) (*GreetHelloResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	return &GreetHelloResponse{
		Msg: fmt.Sprintf("Hello, %s", req.Name),
	}, nil
}
