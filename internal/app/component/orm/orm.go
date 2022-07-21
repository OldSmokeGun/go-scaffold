package orm

import (
	"errors"
	"go-scaffold/internal/app/component/orm/mysql"
	"go-scaffold/internal/app/component/orm/postgres"
	"time"

	klog "github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	zapgorm "moul.io/zapgorm2"
)

var (
	ErrUnsupportedType         = errors.New("unsupported database type")
	ErrUnsupportedResolverType = errors.New("unsupported resolver type")
)

// Driver database driver type
type Driver string

func (d Driver) String() string {
	return string(d)
}

const (
	MySQL       Driver = "mysql"
	PostgresSQL Driver = "postgres"
)

// LogLevel database logger level
type LogLevel string

const (
	Silent LogLevel = "silent"
	Error  LogLevel = "error"
	Warn   LogLevel = "warn"
	Info   LogLevel = "info"
)

// Convert convert to gorm logger level
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

// DSN database connection configuration
type DSN struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Options  []string
}

type Config struct {
	Driver Driver
	DSN
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxIdleTime int64
	ConnMaxLifeTime int64
	LogLevel        LogLevel
	Plugins         func(db *gorm.DB) ([]gorm.Plugin, error)
	Resolvers       []Resolver
}

// New initialize orm
func New(config *Config, kLogger klog.Logger, zLogger *zap.Logger) (db *gorm.DB, cleanup func(), err error) {
	if config == nil {
		return nil, func() {}, nil
	}

	gLogger := zapgorm.New(zLogger.WithOptions(zap.AddCallerSkip(3)))
	gLogger.SetAsDefault()

	switch config.Driver {
	case MySQL:
		db, err = mysql.New(mysql.Config{
			Driver:                    config.Driver.String(),
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
	case PostgresSQL:
		db, err = postgres.New(postgres.Config{
			Driver:               config.Driver.String(),
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

	if len(config.Resolvers) > 0 {
		if err = registerResolver(db, config.Driver, config.Resolvers); err != nil {
			return nil, nil, err
		}
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

func registerResolver(db *gorm.DB, driver Driver, resolvers []Resolver) error {
	if len(resolvers) > 0 {
		var (
			sources  = make([]gorm.Dialector, 0, len(resolvers))
			replicas = make([]gorm.Dialector, 0, len(resolvers))
		)

		for _, resolver := range resolvers {
			dial, err := BuildDialector(driver, resolver.DSN)
			if err != nil {
				return err
			}
			switch resolver.Type {
			case Source:
				sources = append(sources, dial)
			case Replica:
				replicas = append(replicas, dial)
			default:
				return ErrUnsupportedResolverType
			}
		}

		return db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  sources,
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}))
	}

	return nil
}
