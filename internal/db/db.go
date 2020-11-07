package db

import (
	"gin-scaffold/internal/global"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

func Init() error {
	var (
		err         error
		host        = viper.GetString("mysql.host")
		port        = viper.GetString("mysql.port")
		database    = viper.GetString("mysql.database")
		username    = viper.GetString("mysql.username")
		password    = viper.GetString("mysql.password")
		charset     = viper.GetString("mysql.charset")
		maxIdleConn = viper.GetInt("mysql.max_idle_conn")
		maxOpenConn = viper.GetInt("mysql.max_open_conn")
	)

	dialect := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=" + charset + "&parseTime=True&loc=Local"

	global.DB, err = gorm.Open("mysql", dialect)

	if err != nil {
		return err
	}

	global.DB.DB().SetMaxIdleConns(maxIdleConn)
	global.DB.DB().SetMaxOpenConns(maxOpenConn)

	return nil
}
