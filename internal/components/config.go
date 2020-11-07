package components

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

const DefaultConfig = "config.yaml"

func LoadConfig() error {
	var (
		config = DefaultConfig
	)

	if v := flag.Lookup("config").Value.String(); v != "" {
		config = v
	}

	viper.SetConfigName(strings.TrimSuffix(filepath.Base(config), filepath.Ext(config)))
	viper.AddConfigPath(filepath.Dir(config))
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件出错，错误信息：%w", err)
	}

	return nil
}
