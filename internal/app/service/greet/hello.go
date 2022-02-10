package greet

import (
	"context"
	"fmt"
)

type HelloParam struct {
	Name string
}

func (l *service) Hello(ctx context.Context, param *HelloParam) (string, error) {
	return fmt.Sprintf("Hello, %s!", param.Name), nil
}
