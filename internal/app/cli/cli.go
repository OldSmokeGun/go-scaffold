package cli

import (
	"github.com/spf13/cobra"
	"go-scaffold/internal/app/cli/command/greet"
	"go-scaffold/internal/app/cli/command/script"
	"go-scaffold/internal/app/cli/pkg/commandx"
	"go-scaffold/internal/app/global"
)

// Setup 命令行应用初始化入口
func Setup() error {
	cli := commandx.NewCommandLine(global.Command())

	// 注册业务的子命令
	greetHandler := greet.NewHandler()
	cli.RegistryBusiness([]*commandx.Command{
		{
			Entity: &cobra.Command{
				Use:   "greet",
				Short: "示例子命令",
				Run:   greetHandler.Default,
			},
			OptionFunc: func(command *cobra.Command) {
				command.Flags().StringP("example", "e", "foo", "示例 flag")
			},
			Children: []*commandx.Command{
				{
					Entity: &cobra.Command{
						Use:   "to",
						Short: "示例子命令",
						Run:   greetHandler.To,
					},
				},
			},
		},
	})

	// 注册临时脚本命令
	s0000000000 := script.NewS0000000000()
	cli.RegistryScript([]*commandx.Command{
		{
			Entity: &cobra.Command{
				Use:   "S0000000000",
				Short: "示例脚本 S0000000000",
				Run:   s0000000000.Run,
			},
		},
	})

	return nil
}
