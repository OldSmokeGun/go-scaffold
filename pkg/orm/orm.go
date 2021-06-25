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

type Config struct {
	Driver          string
	Host            string
	Port            string
	Database        string
	Username        string
	Password        string
	Options         []string
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifeTime int64
	LogLevel        string
}

// Initialize 初始化 orm
func Initialize(c *Config) (db *gorm.DB, err error) {
	switch c.Driver {
	case "mysql":
		db, err = mysql.GetDB(mysql.Config{
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
			LogLevel:                  LogMode(c.LogLevel).Convert(),
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
		db, err = postgres.GetDB(postgres.Config{
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
			LogLevel:             LogMode(c.LogLevel).Convert(),
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

type LogMode string

func (l LogMode) Convert() logger.Interface {
	switch l {
	case "silent":
		return logger.Default.LogMode(logger.Silent)
	case "error":
		return logger.Default.LogMode(logger.Error)
	case "warn":
		return logger.Default.LogMode(logger.Warn)
	case "info":
		return logger.Default.LogMode(logger.Info)
	default:
		return logger.Default.LogMode(logger.Info)
	}
}
