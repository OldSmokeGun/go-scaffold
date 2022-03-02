package command

import (
	"github.com/google/wire"
	"github.com/spf13/cobra"
	"go-scaffold/internal/app/command/handler"
	"go-scaffold/internal/app/command/handler/greet"
	"go-scaffold/internal/app/command/pkg/commandx"
	"go-scaffold/internal/app/command/script"
	"go-scaffold/internal/app/component"
)

var ProviderSet = wire.NewSet(
	component.ProviderSet,
	script.ProviderSet,
	handler.ProviderSet,
)

func Setup(rootCommand *cobra.Command, newCommand func() (*Command, func(), error)) {
	set := commandx.NewCommandSet(rootCommand)

	// TODO 编写子命令

	// 注册业务的子命令
	set.RegistryBusiness([]*commandx.Command{
		{
			Entity: &cobra.Command{
				Use:   "greet",
				Short: "示例子命令",
				Run: func(cmd *cobra.Command, args []string) {
					command, cleanup, err := newCommand()
					if err != nil {
						panic(err)
					}
					defer cleanup()
					command.greetHandler.Default(cmd, args)
				},
			},
			OptionFunc: func(command *cobra.Command) {
				command.Flags().StringP("example", "e", "foo", "示例 flag")
			},
			Children: []*commandx.Command{
				{
					Entity: &cobra.Command{
						Use:   "to",
						Short: "示例子命令",
						Run: func(cmd *cobra.Command, args []string) {
							command, cleanup, err := newCommand()
							if err != nil {
								panic(err)
							}
							defer cleanup()
							command.greetHandler.To(cmd, args)
						},
					},
				},
			},
		},
	})

	// 注册临时脚本命令
	set.RegistryScript([]*commandx.Command{
		{
			Entity: &cobra.Command{
				Use:   "S0000000000",
				Short: "示例脚本 S0000000000",
				Run: func(cmd *cobra.Command, args []string) {
					command, cleanup, err := newCommand()
					if err != nil {
						panic(err)
					}
					defer cleanup()
					command.s0000000000Script.Run(cmd, args)
				},
			},
		},
	})
}

type Command struct {
	greetHandler      greet.Handler
	s0000000000Script *script.S0000000000
}

func New(
	greetHandler greet.Handler,
	s0000000000Script *script.S0000000000,
) *Command {
	return &Command{
		greetHandler:      greetHandler,
		s0000000000Script: s0000000000Script,
	}
}
