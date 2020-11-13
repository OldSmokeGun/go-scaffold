package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

var Type = "mysql"

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
	ConnMaxLifeTime           time.Duration
	Conn                      gorm.ConnPool
	SkipInitializeWithVersion bool
	DefaultStringSize         uint
	DisableDatetimePrecision  bool
	DontSupportRenameIndex    bool
	DontSupportRenameColumn   bool
	LogLevel                  logger.Interface
}

func (*Config) GetType() string {
	return Type
}

func (c *Config) GetDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName:                c.Driver,
		DSN:                       c.GetDNS(),
		Conn:                      c.Conn,
		SkipInitializeWithVersion: c.SkipInitializeWithVersion,
		DefaultStringSize:         c.DefaultStringSize,
		DisableDatetimePrecision:  c.DisableDatetimePrecision,
		DontSupportRenameIndex:    c.DontSupportRenameIndex,
		DontSupportRenameColumn:   c.DontSupportRenameColumn,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   c.LogLevel,
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

	if c.ConnMaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(c.ConnMaxLifeTime)
	}

	return db, nil
}

func (c *Config) GetDNS() string {
	options := strings.Join(c.Options, "&")
	dsn := c.Username + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.Database + "?" + options
	return dsn
}
