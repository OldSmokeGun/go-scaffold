package main

import (
	"gin-scaffold/app/commands"
	"gin-scaffold/core"
	"gin-scaffold/core/components"
	"gin-scaffold/core/global"
	"gin-scaffold/core/orm"
	"github.com/spf13/cobra"
)

const AppName = "server"

func main() {
	cobra.OnInitialize(boot)

	// 创建根命令
	rootCmd := &cobra.Command{
		Use: AppName,
		Run: func(cmd *cobra.Command, args []string) {
			// 启动内核
			core.Boot()
		},
	}

	// flag 声明与解析
	rootCmd.Flags().StringP("host", "", "", "监听地址")
	rootCmd.Flags().StringP("port", "p", "", "监听端口")
	rootCmd.Flags().StringP("config", "c", "", "配置文件路径")

	// 注册子命令
	commands.Register(rootCmd)

	global.SetRootCommand(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

// 初始化基本依赖
func boot() {
	var flag = global.RootCommand().Flags()

	// 注册配置对象
	if err := components.RegisterConfigurator(flag.Lookup("config").Value.String()); err != nil {
		panic(err)
	}

	// 注册日志对象
	if err := components.RegisterLogger(global.Configurator().GetString("logPath")); err != nil {
		panic(err)
	}

	// 如果获取到数据库配置
	// 则初始化 orm
	if len(global.Configurator().GetStringMap("db")) > 0 {
		of := new(orm.Config)
		if err := global.Configurator().UnmarshalKey("db", of); err != nil {
			panic(err)
		}

		if err := orm.Init(of); err != nil {
			panic(err)
		}

		defer func() {
			sqlDB, err := global.DB().DB()
			if err != nil {
				panic(err)
			}

			if err := sqlDB.Close(); err != nil {
				panic(err)
			}
		}()
	}
}
