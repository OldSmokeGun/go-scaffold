package main

import (
	"gin-scaffold/internal"
	"github.com/spf13/pflag"
)

func main() {
	pflag.StringP("host", "h", "", "监听地址")
	pflag.StringP("port", "p", "", "监听端口")
	pflag.StringP("config", "c", "", "配置文件地址")
	pflag.BoolP("template-from-vfs", "", true, "从虚拟文件系统中使用模版")

	pflag.Parse()

	internal.Bootstrap()
}
