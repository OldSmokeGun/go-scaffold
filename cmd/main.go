package main

import (
	"gin-scaffold/global"
	"gin-scaffold/internal"
	"gin-scaffold/internal/commands"
	appcontext "gin-scaffold/internal/context"
	"gin-scaffold/pkg/configurator"
	"gin-scaffold/pkg/logger"
	"github.com/spf13/cobra"
	"path/filepath"
)

const (
	appName     = "app"       // 命令名称
	defaultHost = "127.0.0.1" // 默认监听地址
	defaultPort = "9527"      // 默认监听端口
)

var (
	appCtx            = appcontext.Default()
	defaultConfigPath = filepath.Join(filepath.Dir(filepath.Dir(global.GetBinPath())), "config/config.yaml") // 默认配置文件路径
	defaultLogPath    = filepath.Join(filepath.Dir(filepath.Dir(global.GetBinPath())), "logs/framework.log") // 默认日志文件路径
)

func main() {
	cobra.OnInitialize(initialize)

	rootCmd := &cobra.Command{
		Use: appName,
		Run: func(cmd *cobra.Command, args []string) {
			internal.Start(cmd, appCtx)
		},
	}

	rootCmd.Flags().StringP("host", "", defaultHost, "监听地址")
	rootCmd.Flags().StringP("port", "p", defaultPort, "监听端口")
	rootCmd.Flags().StringP("config", "c", defaultConfigPath, "配置文件路径")
	rootCmd.Flags().StringP("log", "l", defaultLogPath, "日志文件路径")

	// 子命令注册钩子
	commands.Register(rootCmd)

	// 传递根命令对象
	appCtx.SetRootCommand(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

// 初始化基本依赖
func initialize() {
	var (
		configPath, logPath string
		flag                = appCtx.GetRootCommand().Flags()
	)

	// 注册配置对象
	configPath = flag.Lookup("config").Value.String()
	cfg, err := configurator.Register(configPath)
	if err != nil {
		panic(err)
	}

	// 注册日志对象
	logFlag := flag.Lookup("log")
	if logFlag.Changed {
		logPath = logFlag.Value.String()
	} else {
		if cfg.InConfig("Log") {
			logPath = cfg.GetString("Log")
		} else {
			logPath = logFlag.DefValue
		}
	}
	lg, err := logger.Register(logPath)
	if err != nil {
		panic(err)
	}

	// 设置 conf 对象中的属性
	appCtx.Config.AppConf.Env = cfg.GetString("Env")
	appCtx.Config.AppConf.Log = logPath

	// 传递日志对象
	appCtx.SetLogger(lg)
	// 传递配置对象
	appCtx.SetConfigurator(cfg)
}
