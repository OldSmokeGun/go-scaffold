package main

import (
	"gin-scaffold/app"
	"gin-scaffold/app/commands"
	"gin-scaffold/core"
	"gin-scaffold/core/components"
	"gin-scaffold/core/global"
	"gin-scaffold/core/orm"
	"github.com/spf13/cobra"
	"path/filepath"
)

const (
	appName     = "server"    // 命令名称
	defaultHost = "127.0.0.1" // 默认监听地址
	defaultPort = "9527"      // 默认监听端口
)

var (
	defaultConfigPath = filepath.Join(filepath.Dir(filepath.Dir(global.BinPath())), "config/config.yaml") // 默认配置文件路径
	defaultLogPath    = filepath.Join(filepath.Dir(filepath.Dir(global.BinPath())), "logs/framework.log") // 默认日志文件路径
)

func main() {
	cobra.OnInitialize(boot)

	// 释放资源
	defer func() {
		sqlDB, err := global.DB().DB()
		if err != nil {
			panic(err)
		}

		if err := sqlDB.Close(); err != nil {
			panic(err)
		}
	}()

	// 创建根命令
	rootCmd := &cobra.Command{
		Use: appName,
		Run: func(cmd *cobra.Command, args []string) {
			// 启动内核
			core.Boot()
		},
	}

	// flag 声明与解析
	rootCmd.Flags().StringP("host", "", defaultHost, "监听地址")
	rootCmd.Flags().StringP("port", "p", defaultPort, "监听端口")
	rootCmd.Flags().StringP("config", "c", defaultConfigPath, "配置文件路径")
	rootCmd.Flags().StringP("log", "l", defaultLogPath, "日志文件路径")

	// 注册子命令
	commands.Register(rootCmd)

	global.SetRootCommand(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

// 初始化基本依赖
func boot() {
	var (
		configPath, logPath string
		flag                = global.RootCommand().Flags()
	)

	// 注册配置对象
	configPath = flag.Lookup("config").Value.String()
	if err := components.RegisterConfigurator(configPath); err != nil {
		panic(err)
	}

	// 注册日志对象
	logFlag := flag.Lookup("log")
	if logFlag.Changed {
		logPath = logFlag.Value.String()
	} else {
		if global.Configurator().InConfig("log") {
			logPath = global.Configurator().GetString("log")
		} else {
			logPath = logFlag.DefValue
		}
	}
	if err := components.RegisterLogger(logPath); err != nil {
		panic(err)
	}

	// 如果获取到数据库配置，则初始化 orm
	if len(global.Configurator().GetStringMap("db")) > 0 {
		ormConfig := new(orm.Config)
		if err := global.Configurator().UnmarshalKey("db", ormConfig); err != nil {
			panic(err)
		}

		if err := orm.Init(ormConfig); err != nil {
			panic(err)
		}
	}

	// 框架基本初始化后调用钩子函数
	app.FrameInitialize()
}
