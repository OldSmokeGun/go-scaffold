package core

import (
	"gin-scaffold/app"
	"gin-scaffold/app/routes"
	"gin-scaffold/core/global"
	"github.com/gin-gonic/gin"
)

const (
	// DefaultHost 默认监听 host
	DefaultHost = "127.0.0.1"
	// DefaultPort 默认监听端口
	DefaultPort = "9527"
)

// Boot 引导内核启动
func Boot() {
	var (
		host = DefaultHost
		port = DefaultPort
		r    = gin.Default()
		flag = global.RootCommand().Flags()

		templateGlob string
	)

	if v := global.Configurator().GetString("host"); v != "" {
		host = v
	}

	if v := global.Configurator().GetString("port"); v != "" {
		port = v
	}

	if v := global.Configurator().GetString("templateGlob"); v != "" {
		templateGlob = v
	}

	if v := flag.Lookup("host").Value.String(); v != "" {
		host = v
	}

	if v := flag.Lookup("port").Value.String(); v != "" {
		port = v
	}

	if v := flag.Lookup("template-glob").Value.String(); v != "" {
		templateGlob = v
	}

	if templateGlob != "" {
		r.LoadHTMLGlob(templateGlob)
	}

	// 注册路由
	routes.Register(r)

	// 调用 app 启动前的钩子
	if err := app.Run(r); err != nil {
		panic(err)
	}

	// http 服务启动
	if err := r.Run(host + ":" + port); err != nil {
		panic(err)
	}
}
