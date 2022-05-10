package tests

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	klog "github.com/go-kratos/kratos/v2/log"
	"go-scaffold/pkg/log"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	zapgorm "moul.io/zapgorm2"
)

type DB struct {
	MDB  *sql.DB
	Mock sqlmock.Sqlmock
	DB   *gorm.DB
}

// NewDB 初始化测试 DB
func NewDB(zLogger *zap.Logger, kLogger klog.Logger) (*DB, func(), error) {
	mdb, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gLogger := zapgorm.New(zLogger)
	gLogger.SetAsDefault()

	var logMode glog.LogLevel

	switch log.Level(logLevel) {
	case "silent":
		logMode = glog.Silent
	case log.Debug:
		fallthrough
	case log.Info:
		logMode = glog.Info
	case log.Warn:
		logMode = glog.Warn
	case log.Error:
		fallthrough
	case log.DPanic:
		fallthrough
	case log.Panic:
		fallthrough
	case log.Fatal:
		logMode = glog.Error
	default:
		logMode = glog.Info
	}

	gdb, err := gorm.Open(
		mysql.New(mysql.Config{
			Conn:                      mdb,
			SkipInitializeWithVersion: true,
		}),
		&gorm.Config{
			Logger: gLogger.LogMode(logMode),
		},
	)
	if err != nil {
		return nil, nil, err
	}

	logger := klog.NewHelper(kLogger)

	cleanup := func() {
		sqlDB, err := gdb.DB()
		if err != nil {
			logger.Error(err)
		}
		if err := sqlDB.Close(); err != nil {
			logger.Error(err)
		}

		if err := mdb.Close(); err != nil {
			logger.Error(err)
		}
	}

	return &DB{
		MDB:  mdb,
		Mock: mock,
		DB:   gdb,
	}, cleanup, nil
}
