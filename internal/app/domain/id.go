package domain

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

type ID int64

func (i ID) ValidateWithContext(ctx context.Context) error {
	scene := GetSceneFromContext(ctx)
	return errors.WithStack(validation.Validate(int64(i),
		validation.When(
			scene == Detail || scene == Update || scene == Delete,
			validation.Required.Error("id 不能为空"),
		)))
}

func (i ID) Int64() int64 {
	return int64(i)
}
