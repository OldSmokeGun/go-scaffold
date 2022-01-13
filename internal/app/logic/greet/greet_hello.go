package greet

import (
	"fmt"
)

type HelloParam struct {
	Name string
}

func (l *logic) Hello(param *HelloParam) (string, error) {
	return fmt.Sprintf("Hello, %s!", param.Name), nil
}
