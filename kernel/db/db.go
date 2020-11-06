package db

import (
	"gin-scaffold/kernel/global"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

func Init() error {
	var (
		err      error
		host     = viper.GetString("mysql.host")
		port     = viper.GetString("mysql.port")
		database = viper.GetString("mysql.database")
		username = viper.GetString("mysql.username")
		password = viper.GetString("mysql.password")
		charset  = viper.GetString("mysql.charset")
	)

	dialect := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=" + charset + "&parseTime=True&loc=Local"

	global.DB, err = gorm.Open("mysql", dialect)

	if err != nil {
		return err
	}

	global.DB.DB().SetMaxIdleConns(20)
	global.DB.DB().SetMaxOpenConns(50)

	return nil
}
