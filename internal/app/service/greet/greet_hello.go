package greet

import (
	"fmt"
)

type HelloParam struct {
	Name string
}

func (l *service) Hello(param *HelloParam) (string, error) {
	return fmt.Sprintf("Hello, %s!", param.Name), nil
}
