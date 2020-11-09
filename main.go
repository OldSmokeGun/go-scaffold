package main

import (
	"gin-scaffold/internal"
	"github.com/spf13/pflag"
)

func main() {
	pflag.StringP("host", "h", "", "监听地址")
	pflag.StringP("port", "p", "", "监听端口")
	pflag.StringP("config", "c", "", "配置文件地址")

	pflag.Parse()

	internal.Bootstrap()
}
