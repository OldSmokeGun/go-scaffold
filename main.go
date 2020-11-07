package main

import (
	"flag"
	"gin-scaffold/internal"
	"gin-scaffold/internal/components"
	"gin-scaffold/internal/db"
	"gin-scaffold/internal/global"
)

func main() {
	flag.String("config", "", "配置文件地址")
	flag.String("host", "", "监听地址")
	flag.String("port", "", "监听端口")

	flag.Parse()

	if err := components.LoadConfig(); err != nil {
		panic(err)
	}

	if err := components.InitLogger(); err != nil {
		panic(err)
	}

	if err := db.Init(); err != nil {
		panic(err)
	}

	defer func() {
		if err := global.DB.Close(); err != nil {
			panic(err)
		}
	}()

	internal.Bootstrap()
}
