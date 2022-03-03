package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

type Config struct {
	Driver                    string
	Host                      string
	Port                      string
	Database                  string
	Username                  string
	Password                  string
	Options                   []string
	MaxIdleConn               int
	MaxOpenConn               int
	ConnMaxIdleTime           time.Duration
	ConnMaxLifeTime           time.Duration
	Logger                    logger.Interface
	Conn                      gorm.ConnPool
	SkipInitializeWithVersion bool
	DefaultStringSize         uint
	DisableDatetimePrecision  bool
	DontSupportRenameIndex    bool
	DontSupportRenameColumn   bool
}

// New 返回 *gorm.DB
func New(c Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName:                c.Driver,
		DSN:                       buildDNS(c),
		Conn:                      c.Conn,
		SkipInitializeWithVersion: c.SkipInitializeWithVersion,
		DefaultStringSize:         c.DefaultStringSize,
		DisableDatetimePrecision:  c.DisableDatetimePrecision,
		DontSupportRenameIndex:    c.DontSupportRenameIndex,
		DontSupportRenameColumn:   c.DontSupportRenameColumn,
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
		sqlDB.SetConnMaxLifetime(c.ConnMaxIdleTime)
	}

	if c.ConnMaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(c.ConnMaxLifeTime)
	}

	return db, nil
}

// buildDNS 构建连接数据库的 dns
func buildDNS(c Config) string {
	options := strings.Join(c.Options, "&")
	dsn := c.Username + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.Database + "?" + options
	return dsn
}
