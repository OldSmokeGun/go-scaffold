package main

import (
	"gin-scaffold/global"
	"gin-scaffold/internal"
	"gin-scaffold/internal/commands"
	"gin-scaffold/internal/config"
	"gin-scaffold/internal/ctx"
	"gin-scaffold/pkg/configurator"
	"gin-scaffold/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"path/filepath"
	"strconv"
)

const (
	appName     = "app"       // 命令名称
	defaultHost = "127.0.0.1" // 默认监听地址
	defaultPort = 9527        // 默认监听端口
)

var (
	appCtx            = ctx.Default()
	defaultConfigPath = filepath.Join(filepath.Dir(filepath.Dir(global.GetBinPath())), "config/config.yaml") // 默认配置文件路径
	defaultLogPath    = filepath.Join(filepath.Dir(filepath.Dir(global.GetBinPath())), "logs/framework.log") // 默认日志文件路径
)

func main() {
	// 初始化基本参数
	initialize()

	rootCmd := &cobra.Command{
		Use: appName,
		Run: func(cmd *cobra.Command, args []string) {
			internal.Start(cmd, appCtx)
		},
	}

	// 子命令注册钩子
	commands.Register(rootCmd, appCtx)

	// 传递根命令对象
	appCtx.SetRootCommand(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

// initialize 初始化基本参数
func initialize() {
	var (
		configPath string
		logPath    string
		host       string
		port       int
		env        string
	)

	pflag.StringVarP(&configPath, "config", "c", defaultConfigPath, "配置文件路径")
	pflag.StringP("log", "l", defaultLogPath, "日志文件路径")
	pflag.StringP("host", "", defaultHost, "监听地址")
	pflag.IntP("port", "p", defaultPort, "监听端口")
	pflag.StringP("env", "", defaultHost, "运行环境")

	pflag.Parse()

	// 构建配置对象
	cfg, err := configurator.Build(configPath)
	if err != nil {
		panic(err)
	}

	// 构建日志对象
	logFlag := pflag.Lookup("log")
	if logFlag.Changed {
		logPath = logFlag.Value.String()
	} else {
		if cfg.InConfig("Log") {
			logPath = cfg.GetString("Log")
		} else {
			logPath = logFlag.DefValue
		}
	}
	log, err := logger.Build(logPath)
	if err != nil {
		panic(err)
	}

	hostFlag := pflag.Lookup("host")
	if hostFlag.Changed {
		host = hostFlag.Value.String()
	} else {
		if cfg.InConfig("Host") {
			host = cfg.GetString("Host")
		} else {
			host = hostFlag.DefValue
		}
	}

	portFlag := pflag.Lookup("port")
	if portFlag.Changed {
		port, err = strconv.Atoi(portFlag.Value.String())
		if err != nil {
			panic(err)
		}
	} else {
		if cfg.InConfig("Port") {
			port = cfg.GetInt("Port")
		} else {
			port, err = strconv.Atoi(portFlag.DefValue)
			if err != nil {
				panic(err)
			}
		}
	}

	envFlag := pflag.Lookup("env")
	if envFlag.Changed {
		env = envFlag.Value.String()
	} else {
		if cfg.InConfig("Env") {
			env = cfg.GetString("Env")
		} else {
			env = envFlag.DefValue
		}
	}

	appCtx.Config.AppConf = config.AppConf{
		Host: host,
		Port: port,
		Env:  env,
		Log:  logPath,
	}
	// 传递配置对象
	appCtx.SetConfigurator(cfg)
	// 传递日志对象
	appCtx.SetLogger(log)

	// 应用全局初始化
	internal.Initialize(appCtx)
}
