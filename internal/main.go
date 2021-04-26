package main

import (
	"gin-scaffold/app"
	"gin-scaffold/app/commands"
	"gin-scaffold/internal/components/configurator"
	"gin-scaffold/internal/components/logger"
	"gin-scaffold/internal/config"
	"gin-scaffold/internal/global"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"path/filepath"
	"strconv"
)

const (
	appName     = "server"    // 命令名称
	defaultHost = "127.0.0.1" // 默认监听地址
	defaultPort = "9527"      // 默认监听端口
)

var (
	conf              = config.Config{}
	defaultConfigPath = filepath.Join(filepath.Dir(filepath.Dir(global.GetBinPath())), "config/config.yaml") // 默认配置文件路径
	defaultLogPath    = filepath.Join(filepath.Dir(filepath.Dir(global.GetBinPath())), "logs/framework.log") // 默认日志文件路径
)

func main() {
	cobra.OnInitialize(initialize)

	// 创建根命令
	rootCmd := &cobra.Command{
		Use: appName,
		Run: func(cmd *cobra.Command, args []string) {
			Start(cmd, conf)
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

// Start 启动服务
func Start(cmd *cobra.Command, conf config.Config) {
	var (
		err  error
		host string
		port int
		flag = cmd.Flags()
	)

	hostFlag := flag.Lookup("host")
	portFlag := flag.Lookup("port")

	if hostFlag.Changed {
		host = hostFlag.Value.String()
	} else {
		if global.GetConfigurator().InConfig("Host") {
			host = global.GetConfigurator().GetString("Host")
		} else {
			host = hostFlag.DefValue
		}
	}

	if portFlag.Changed {
		port, err = strconv.Atoi(portFlag.Value.String())
		if err != nil {
			panic(err)
		}
	} else {
		if global.GetConfigurator().InConfig("Port") {
			port = global.GetConfigurator().GetInt("Port")
		} else {
			port, err = strconv.Atoi(portFlag.DefValue)
			if err != nil {
				panic(err)
			}
		}
	}

	// 设置 conf 对象中的属性
	conf.Host = host
	conf.Port = port

	// 启动 app
	app.Start(gin.Default(), conf)
}

// 初始化基本依赖
func initialize() {
	var (
		configPath, logPath string
		flag                = global.GetRootCommand().Flags()
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
	conf.Env = cfg.GetString("Env")
	conf.Log = logPath

	// 设置全局日志对象
	global.SetLogger(lg)
	// 设置全局配置对象
	global.SetConfigurator(cfg)

	// 框架基本初始化后调用钩子函数
	if err := app.Initialize(); err != nil {
		panic(err)
	}
}
