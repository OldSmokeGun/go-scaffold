package main

import (
	"gin-scaffold/core"
	"gin-scaffold/core/components"
	"gin-scaffold/core/global"
	"gin-scaffold/core/orm"
	"github.com/spf13/pflag"
)

func main() {
	// flag 声明与解析
	pflag.StringP("host", "h", "", "监听地址")
	pflag.StringP("port", "p", "", "监听端口")
	pflag.StringP("config", "c", "", "配置文件路径")
	pflag.StringP("template-glob", "t", "", "模板文件 glob 表达式")

	pflag.Parse()

	// 注册配置对象
	if err := components.RegisterConfigurator(pflag.Lookup("config").Value.String()); err != nil {
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

	// 启动内核
	core.Bootstrap()
}
