package kernel

import (
	"flag"
	"gin-scaffold/kernel/routers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strings"
)

const DefaultHost = "0.0.0.0"
const DefaultPort = "9527"

func Bootstrap() {
	var (
		host string
		port string
	)

	flag.StringVar(&host, "host", "", "监听地址")
	flag.StringVar(&port, "port", "", "监听端口")
	flag.Parse()

	if host == "" {
		host = DefaultHost
		if viper.InConfig("host") {
			host = viper.GetString("host")
		}
	}

	if port == "" {
		port = DefaultPort
		if viper.InConfig("port") {
			port = viper.GetString("port")
		}
	}

	r := gin.Default()
	routers.Register(r)

	listenAddr := strings.Join([]string{host, port}, ":")

	if err := r.Run(listenAddr); err != nil {
		panic(err)
	}
}
