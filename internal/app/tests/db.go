package tests

import (
	"database/sql"
	"time"

	igorm "go-scaffold/internal/app/pkg/gorm"
	"go-scaffold/internal/config"
	glog "go-scaffold/pkg/log/gorm"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
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

	dialect, err := igorm.NewDialect(config.DBDriver(dbDriver), mdb)
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
			logger.Error("get sql db error", err)
		}
		if err := sdb.Close(); err != nil {
			logger.Error("close sql db error", err)
		}

		if err := mdb.Close(); err != nil {
			logger.Error("close mock db error", err)
		}
	}

	return &DB{
		MDB:  mdb,
		Mock: mock,
		DB:   gdb,
	}, cleanup, nil
}
