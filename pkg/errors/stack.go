package errors

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// StackTrace handle the error carrying the stack
type StackTrace errors.StackTrace

func (st StackTrace) MarshalJSON() ([]byte, error) {
	var stacks []string
	for _, frame := range st {
		f, err := frame.MarshalText()
		if err != nil {
			return nil, err
		}
		stacks = append(stacks, string(f))
	}
	return json.Marshal(stacks)
}

func (st StackTrace) Format(s fmt.State, verb rune) {
	io.WriteString(s, fmt.Sprintf("%+v", st))
}

// ErrorStackTrace format error with stack
func ErrorStackTrace(err error) StackTrace {
	if v, ok := err.(interface {
		StackTrace() errors.StackTrace
	}); ok {
		return StackTrace(v.StackTrace())
	}
	return nil
}

// IsStackTrace check if error implements the stack
func IsStackTrace(err error) bool {
	_, ok := err.(interface {
		StackTrace() errors.StackTrace
	})
	return ok
}
