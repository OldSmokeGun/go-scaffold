package main

import (
	"flag"
	"gin-scaffold/internal"
)

func main() {
	flag.String("config", "", "配置文件地址")
	flag.String("host", "", "监听地址")
	flag.String("port", "", "监听端口")

	flag.Parse()

	internal.Bootstrap()
}
