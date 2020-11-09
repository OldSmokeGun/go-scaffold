package internal

import (
	"errors"
	"gin-scaffold/internal/components"
	"gin-scaffold/internal/db"
	"gin-scaffold/internal/db/mysql"
	"gin-scaffold/internal/db/postgres"
	"gin-scaffold/internal/global"
	"gin-scaffold/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
	"time"
)

const DefaultHost = "127.0.0.1"
const DefaultPort = "9527"

func Bootstrap() {
	var (
		err  error
		host = DefaultHost
		port = DefaultPort
	)

	global.BinPath, err = os.Executable()
	if err != nil {
		panic(err)
	}

	if err = components.LoadViper(pflag.Lookup("config").Value.String()); err != nil {
		panic(err)
	}

	if err = components.LoadLogrus(viper.GetString("errors_log")); err != nil {
		panic(err)
	}

	global.DB, err = initDB()
	if err != nil {
		panic(err)
	}
	defer func() {
		sqlDB, err := global.DB.DB()
		if err != nil {
			panic(err)
		}

		err = sqlDB.Close()
		if err != nil {
			panic(err)
		}
	}()

	if v := pflag.Lookup("host").Value.String(); v != "" {
		host = v
	} else {
		if viper.GetString("host") != "" {
			host = viper.GetString("host")
		}
	}

	if v := pflag.Lookup("port").Value.String(); v != "" {
		port = v
	} else {
		if viper.GetString("port") != "" {
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

func initDB() (*gorm.DB, error) {
	var (
		config db.Config
		dbType = viper.GetString("db.type")
	)

	if len(viper.GetStringMap("db")) == 0 {
		return nil, nil
	}

	switch dbType {
	case "mysql":
		config = &mysql.Config{
			Driver:          viper.GetString("db.driver"),
			Host:            viper.GetString("db.host"),
			Port:            viper.GetString("db.port"),
			Database:        viper.GetString("db.database"),
			Username:        viper.GetString("db.username"),
			Password:        viper.GetString("db.password"),
			Options:         viper.GetStringSlice("db.options"),
			MaxIdleConn:     viper.GetInt("db.max_idle_conn"),
			MaxOpenConn:     viper.GetInt("db.max_open_conn"),
			ConnMaxLifeTime: time.Second * viper.GetDuration("db.conn_max_life_time"),
		}
	case "postgres":
		config = &postgres.Config{
			Driver:          viper.GetString("db.driver"),
			Host:            viper.GetString("db.host"),
			Port:            viper.GetString("db.port"),
			Database:        viper.GetString("db.database"),
			Username:        viper.GetString("db.username"),
			Password:        viper.GetString("db.password"),
			Options:         viper.GetStringSlice("db.options"),
			MaxIdleConn:     viper.GetInt("db.max_idle_conn"),
			MaxOpenConn:     viper.GetInt("db.max_open_conn"),
			ConnMaxLifeTime: time.Second * viper.GetDuration("db.conn_max_life_time"),
		}
	default:
		return nil, errors.New("unsupported database type " + dbType)
	}

	DB, err := db.Init(config)
	if err != nil {
		return nil, err
	}

	return DB, nil
}
