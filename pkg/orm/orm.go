package orm

import (
	"errors"
	"gin-scaffold/pkg/orm/mysql"
	"gin-scaffold/pkg/orm/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var (
	ErrUnsupportedType = errors.New("unsupported database type")
)

// New 初始化 orm，返回 *gorm.DB
func New(c Config) (db *gorm.DB, err error) {
	switch c.Driver {
	case "mysql":
		db, err = mysql.New(mysql.Config{
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
			Logger:                    logger.Default.LogMode(c.LogLevel.Convert()),
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
		db, err = postgres.New(postgres.Config{
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
			Logger:               logger.Default.LogMode(c.LogLevel.Convert()),
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

// MustNew 初始化 orm，返回 *gorm.DB
func MustNew(c Config) *gorm.DB {
	db, err := New(c)
	if err != nil {
		panic(err)
	}

	return db
}
