package postgres

import (
	"database/sql"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Driver               string
	Addr                 string
	Database             string
	Username             string
	Password             string
	Options              string
	MaxIdleConn          int
	MaxOpenConn          int
	ConnMaxIdleTime      time.Duration
	ConnMaxLifeTime      time.Duration
	Logger               logger.Interface
	Conn                 *sql.DB
	PreferSimpleProtocol bool
}

// New initialize *gorm.DB
func New(c Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName:           c.Driver,
		DSN:                  BuildDSN(c),
		PreferSimpleProtocol: c.PreferSimpleProtocol,
		Conn:                 c.Conn,
	}), &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   c.Logger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if c.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	}

	if c.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	}

	if c.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(c.ConnMaxIdleTime)
	}

	if c.ConnMaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(c.ConnMaxLifeTime)
	}

	return db, nil
}

// BuildDSN build dss to connect to the database
func BuildDSN(c Config) string {
	var host, port string
	s := strings.SplitN(c.Addr, ":", 2)
	if len(s) == 2 {
		host = s[0]
		port = s[1]
	}

	dsn := "host=" + host + " port=" + port + " user=" + c.Username + " password=" + c.Password + " dbname=" + c.Database + " " + c.Options
	return dsn
}
