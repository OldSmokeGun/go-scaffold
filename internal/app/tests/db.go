package tests

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	igorm "go-scaffold/internal/app/pkg/gorm"
	"go-scaffold/internal/config"
	glog "go-scaffold/pkg/log/gorm"
)

type DB struct {
	MDB  *sql.DB
	Mock sqlmock.Sqlmock
	DB   *gorm.DB
}

// NewDB init mock db
func NewDB(logger *slog.Logger) (*DB, func(), error) {
	mdb, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormlogger.Default = glog.NewLogger(logger, glog.Config{
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: false,
		LogInfo:                   true,
	})

	dialect, err := igorm.NewDialect(config.DatabaseDriver(dbDriver), mdb)
	if err != nil {
		return nil, nil, err
	}
	gdb, err := gorm.Open(dialect, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		sdb, err := gdb.DB()
		if err != nil {
			logger.Error("get sql db error", slog.Any("error", err))
		}
		if err := sdb.Close(); err != nil {
			logger.Error("close sql db error", slog.Any("error", err))
		}

		if err := mdb.Close(); err != nil {
			logger.Error("close mock db error", slog.Any("error", err))
		}
	}

	return &DB{
		MDB:  mdb,
		Mock: mock,
		DB:   gdb,
	}, cleanup, nil
}
