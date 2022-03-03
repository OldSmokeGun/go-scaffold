package orm

import (
	"errors"
	klog "github.com/go-kratos/kratos/v2/log"
	"go-scaffold/internal/app/component/orm/mysql"
	"go-scaffold/internal/app/component/orm/postgres"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	zapgorm "moul.io/zapgorm2"
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
	ConnMaxIdleTime int64
	ConnMaxLifeTime int64
	LogLevel        LogLevel
}

type LogLevel string

const (
	Silent LogLevel = "silent"
	Error  LogLevel = "error"
	Warn   LogLevel = "warn"
	Info   LogLevel = "info"
)

func (l LogLevel) Convert() logger.LogLevel {
	switch l {
	case Silent:
		return logger.Silent
	case Error:
		return logger.Error
	case Warn:
		return logger.Warn
	case Info:
		return logger.Info
	default:
		return logger.Info
	}
}

// New 初始化 orm
func New(config *Config, kLogger klog.Logger, zLogger *zap.Logger) (db *gorm.DB, cleanup func(), err error) {
	if config == nil {
		return nil, func() {}, nil
	}

	gLogger := zapgorm.New(zLogger.WithOptions(zap.AddCallerSkip(3)))
	gLogger.SetAsDefault()

	switch config.Driver {
	case "mysql":
		db, err = mysql.New(mysql.Config{
			Driver:                    config.Driver,
			Host:                      config.Host,
			Port:                      config.Port,
			Database:                  config.Database,
			Username:                  config.Username,
			Password:                  config.Password,
			Options:                   config.Options,
			MaxIdleConn:               config.MaxIdleConn,
			MaxOpenConn:               config.MaxOpenConn,
			ConnMaxIdleTime:           time.Second * time.Duration(config.ConnMaxIdleTime),
			ConnMaxLifeTime:           time.Second * time.Duration(config.ConnMaxLifeTime),
			Logger:                    gLogger.LogMode(config.LogLevel.Convert()),
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
			Driver:               config.Driver,
			Host:                 config.Host,
			Port:                 config.Port,
			Database:             config.Database,
			Username:             config.Username,
			Password:             config.Password,
			Options:              config.Options,
			MaxIdleConn:          config.MaxIdleConn,
			MaxOpenConn:          config.MaxOpenConn,
			ConnMaxIdleTime:      time.Second * time.Duration(config.ConnMaxIdleTime),
			ConnMaxLifeTime:      time.Second * time.Duration(config.ConnMaxLifeTime),
			Logger:               gLogger.LogMode(config.LogLevel.Convert()),
			Conn:                 nil,
			PreferSimpleProtocol: false,
		})
		if err != nil {
			return
		}
	default:
		return nil, nil, ErrUnsupportedType
	}

	cleanup = func() {
		klog.NewHelper(kLogger).Info("closing the database resources")

		sqlDB, err := db.DB()
		if err != nil {
			klog.NewHelper(kLogger).Error(err)
		}

		if err := sqlDB.Close(); err != nil {
			klog.NewHelper(kLogger).Error(err)
		}
	}

	return db, cleanup, nil
}
