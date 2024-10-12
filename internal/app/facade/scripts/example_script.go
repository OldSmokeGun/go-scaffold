package scripts

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"go-scaffold/internal/app/controller"
)

type ExampleCmd struct {
	controller *controller.GreetController
}

func NewExampleCmd(
	controller *controller.GreetController,
) *ExampleCmd {
	return &ExampleCmd{
		controller: controller,
	}
}

func (c *ExampleCmd) Run(cmd *cobra.Command) error {
	ret, err := c.controller.Hello(cmd.Context(), controller.GreetHelloRequest{Name: "Example"})
	if err != nil {
		return err
	}

	s, err := jsoniter.MarshalToString(ret)
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println(s)

	return nil
}
