package internal

import (
	"flag"
	"gin-scaffold/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const DefaultHost = "127.0.0.1"
const DefaultPort = "9527"

func Bootstrap() {
	var (
		host = DefaultHost
		port = DefaultPort
	)

	if v := flag.Lookup("host").Value.String(); v != "" {
		host = v
	} else {
		if viper.IsSet("host") && viper.GetString("host") != "" {
			host = viper.GetString("host")
		}
	}

	if v := flag.Lookup("port").Value.String(); v != "" {
		port = v
	} else {
		if viper.IsSet("port") && viper.GetString("port") != "" {
			port = viper.GetString("port")
		}
	}

	r := gin.Default()
	router.Register(r)

	listenAddr := host + ":" + port

	if err := r.Run(listenAddr); err != nil {
		panic(err)
	}
}
