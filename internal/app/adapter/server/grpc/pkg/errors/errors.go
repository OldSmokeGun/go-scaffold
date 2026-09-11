package errors

import (
	"errors"

	kerr "github.com/go-kratos/kratos/v2/errors"

	berr "go-scaffold/internal/errors"
)

// Wrap application internal error
func Wrap(err error) error {
	var se *berr.Error
	if errors.As(err, &se) {
		return kerr.New(se.Code(), se.Reason(), se.Msg())
	}

	e := berr.ErrInternalError
	return kerr.New(e.Code(), e.Reason(), err.Error())
}
