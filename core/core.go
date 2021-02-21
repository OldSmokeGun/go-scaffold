package core

import (
	"gin-scaffold/app"
	"gin-scaffold/app/routes"
	"gin-scaffold/core/global"
	"github.com/gin-gonic/gin"
)

// Boot 引导内核启动
func Boot() {
	var (
		host, port string
		router     = gin.Default()
		flag       = global.RootCommand().Flags()
	)

	hostFlag := flag.Lookup("host")
	portFlag := flag.Lookup("port")

	if hostFlag.Changed {
		host = hostFlag.Value.String()
	} else {
		if global.Configurator().InConfig("host") {
			host = global.Configurator().GetString("host")
		} else {
			host = hostFlag.DefValue
		}
	}

	if portFlag.Changed {
		port = portFlag.Value.String()
	} else {
		if global.Configurator().InConfig("port") {
			port = global.Configurator().GetString("port")
		} else {
			port = portFlag.DefValue
		}
	}

	// 注册路由
	routes.Register(router)

	// 调用 app 启动前的钩子
	if err := app.Initialize(router); err != nil {
		panic(err)
	}

	// http 服务启动
	if err := router.Run(host + ":" + port); err != nil {
		panic(err)
	}
}
