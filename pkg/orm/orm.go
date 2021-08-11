package orm

import (
	"errors"
	"gin-scaffold/pkg/orm/mysql"
	"gin-scaffold/pkg/orm/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"
)

var (
	ErrUnsupportedType = errors.New("unsupported database type")
)

// Setup 初始化 orm，返回 *gorm.DB
func Setup(c Config) (db *gorm.DB, err error) {
	// 设置 logger
	l := logger.New(log.New(io.MultiWriter(c.Output, os.Stdout), "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      c.LogLevel.Convert(),
		Colorful:      false,
	})

	switch c.Driver {
	case "mysql":
		db, err = mysql.NewDB(mysql.Config{
			Driver:                    c.Driver,
			Host:                      c.Host,
			Port:                      c.Port,
			Database:                  c.Database,
			Username:                  c.Username,
			Password:                  c.Password,
			Options:                   c.Options,
			MaxIdleConn:               c.MaxIdleConn,
			MaxOpenConn:               c.MaxOpenConn,
			ConnMaxLifeTime:           time.Second * time.Duration(c.ConnMaxLifeTime),
			Logger:                    l,
			Conn:                      nil,
			SkipInitializeWithVersion: false,
			DefaultStringSize:         0,
			DisableDatetimePrecision:  false,
			DontSupportRenameIndex:    false,
			DontSupportRenameColumn:   false,
		})
		if err != nil {
			return
		}
	case "postgres":
		db, err = postgres.NewDB(postgres.Config{
			Driver:               c.Driver,
			Host:                 c.Host,
			Port:                 c.Port,
			Database:             c.Database,
			Username:             c.Username,
			Password:             c.Password,
			Options:              c.Options,
			MaxIdleConn:          c.MaxIdleConn,
			MaxOpenConn:          c.MaxOpenConn,
			ConnMaxLifeTime:      time.Second * time.Duration(c.ConnMaxLifeTime),
			Logger:               l,
			Conn:                 nil,
			PreferSimpleProtocol: false,
		})
		if err != nil {
			return
		}
	default:
		return nil, ErrUnsupportedType
	}

	return db, nil
}

// MustSetup 初始化 orm，返回 *gorm.DB
func MustSetup(c Config) *gorm.DB {
	db, err := Setup(c)
	if err != nil {
		panic(err)
	}

	return db
}
