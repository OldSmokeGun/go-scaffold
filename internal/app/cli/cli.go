package cli

import (
	"github.com/spf13/cobra"
	"go-scaffold/internal/app/cli/command/greet"
	"go-scaffold/internal/app/cli/pkg/commandx"
	"go-scaffold/internal/app/cli/script"
	"go-scaffold/internal/app/global"
)

// Setup 命令行应用初始化
func Setup() error {
	// 根命令
	rootCommand := global.Command()

	cli := commandx.NewCommandLine(rootCommand)

	// TODO 编写子命令

	// 注册业务的子命令
	cli.RegistryBusiness([]*commandx.Command{
		{
			Entity: &cobra.Command{
				Use:   "greet",
				Short: "示例子命令",
				Run: func(cmd *cobra.Command, args []string) {
					greetHandler := greet.NewHandler()
					greetHandler.Default(cmd, args)
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
							greetHandler := greet.NewHandler()
							greetHandler.To(cmd, args)
						},
					},
				},
			},
		},
	})

	// 注册临时脚本命令
	cli.RegistryScript([]*commandx.Command{
		{
			Entity: &cobra.Command{
				Use:   "S0000000000",
				Short: "示例脚本 S0000000000",
				Run: func(cmd *cobra.Command, args []string) {
					s0000000000 := script.NewS0000000000()
					s0000000000.Run(cmd, args)
				},
			},
		},
	})

	return nil
}
