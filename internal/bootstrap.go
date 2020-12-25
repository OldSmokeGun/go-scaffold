package internal

import (
	"errors"
	"gin-scaffold/internal/app"
	"gin-scaffold/internal/components"
	"gin-scaffold/internal/components/vfs"
	"gin-scaffold/internal/db"
	"gin-scaffold/internal/db/mysql"
	"gin-scaffold/internal/db/postgres"
	"gin-scaffold/internal/global"
	"gin-scaffold/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultHost         = "127.0.0.1"
	DefaultPort         = "9527"
	DefaultTemplateGlob = "*"
)

func Bootstrap() {
	var (
		err          error
		host         = DefaultHost
		port         = DefaultPort
		templateGlob = DefaultTemplateGlob
		appPath      = filepath.Dir(global.BinPath) + "/internal/app"
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

	if d := viper.GetString("templates_glob"); d != "" {
		templateGlob = d
	}

	r := gin.Default()

	if v := pflag.Lookup("template-from-vfs").Value.String(); v == "true" {
		t, err := vfs.LoadTemplatesFromFilesystem()
		if err != nil {
			panic(err)
		}
		r.SetHTMLTemplate(t)
	} else {
		glob := appPath + "/templates/" + templateGlob

		matches, err := filepath.Glob(glob)
		if err != nil {
			panic(err)
		}

		if len(matches) > 0 {
			r.LoadHTMLGlob(glob)
		}
	}

	router.Register(r)

	if err := app.Constructor(r); err != nil {
		panic(err)
	}

	listenAddr := host + ":" + port

	if err := r.Run(listenAddr); err != nil {
		panic(err)
	}
}

func initDB() (*gorm.DB, error) {
	var (
		config     db.Config
		dbType     = viper.GetString("db.type")
		dbLogLevel = logger.Default.LogMode(logger.Info)
	)

	if len(viper.GetStringMap("db")) == 0 {
		return nil, nil
	}

	if viper.GetString("db.log_level") != "" {
		switch viper.GetString("db.log_level") {
		case "silent":
			dbLogLevel = logger.Default.LogMode(logger.Silent)
		case "error":
			dbLogLevel = logger.Default.LogMode(logger.Error)
		case "warn":
			dbLogLevel = logger.Default.LogMode(logger.Warn)
		case "info":
			dbLogLevel = logger.Default.LogMode(logger.Info)
		default:
			dbLogLevel = logger.Default.LogMode(logger.Info)
		}
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
			LogLevel:        dbLogLevel,
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
			LogLevel:        dbLogLevel,
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
