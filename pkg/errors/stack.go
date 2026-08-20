package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"

	"github.com/pkg/errors"
)

// StackTracer is implemented by errors that carry a stack trace.
type StackTracer interface {
	StackTrace() errors.StackTrace
}

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
	if v, ok := err.(StackTracer); ok {
		return StackTrace(v.StackTrace())
	}
	return nil
}

// IsStackTrace check if error implements the stack
func IsStackTrace(err error) bool {
	_, ok := err.(StackTracer)
	return ok
}

// WithStack annotates err with a stack when it has none.
func WithStack(err error, skip int) error {
	if err == nil {
		return nil
	}
	if IsStackTrace(err) {
		return err
	}
	return &withStack{
		error: err,
		stack: callers(3 + skip),
	}
}

type withStack struct {
	error
	stack []uintptr
}

func (w *withStack) Unwrap() error { return w.error }

func (w *withStack) StackTrace() errors.StackTrace {
	f := make([]errors.Frame, len(w.stack))
	for i := range w.stack {
		f[i] = errors.Frame(w.stack[i])
	}
	return f
}

func callers(skip int) []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	st := make([]uintptr, n)
	copy(st, pcs[:n])
	return st
}
