package main

import (
	"gin-scaffold/core"
	"gin-scaffold/core/components"
	"gin-scaffold/core/global"
	"gin-scaffold/core/orm"
	"github.com/spf13/pflag"
	"time"
)

func main() {
	// flag 声明与解析
	pflag.StringP("host", "h", "", "监听地址")
	pflag.StringP("port", "p", "", "监听端口")
	pflag.StringP("config", "c", "", "配置文件路径")
	pflag.StringP("templates_glob", "t", "", "模板文件 glob 表达式")

	pflag.Parse()

	// 注册配置对象
	if err := components.RegisterConfigurator(pflag.Lookup("config").Value.String()); err != nil {
		panic(err)
	}

	// 注册日志对象
	if err := components.RegisterLogger(global.Configurator().GetString("log_path")); err != nil {
		panic(err)
	}

	// 如果获取到数据库配置
	// 则初始化 orm
	if len(global.Configurator().GetStringMap("db")) > 0 {
		if err := orm.Init(orm.Config{
			Driver:          global.Configurator().GetString("db.driver"),
			Host:            global.Configurator().GetString("db.host"),
			Port:            global.Configurator().GetString("db.port"),
			Database:        global.Configurator().GetString("db.database"),
			Username:        global.Configurator().GetString("db.username"),
			Password:        global.Configurator().GetString("db.password"),
			Options:         global.Configurator().GetStringSlice("db.options"),
			MaxIdleConn:     global.Configurator().GetInt("db.max_idle_conn"),
			MaxOpenConn:     global.Configurator().GetInt("db.max_open_conn"),
			ConnMaxLifeTime: time.Second * global.Configurator().GetDuration("db.conn_max_life_time"),
			LogLevel:        orm.LogMode(global.Configurator().GetString("db.log_level")).Convert(),
		}); err != nil {
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
