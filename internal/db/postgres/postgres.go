package postgres

import (
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"time"
)

var Type = "postgres"

type Config struct {
	Driver               string
	Host                 string
	Port                 string
	Database             string
	Username             string
	Password             string
	Options              []string
	MaxIdleConn          int
	MaxOpenConn          int
	ConnMaxLifeTime      time.Duration
	PreferSimpleProtocol bool
	Conn                 *sql.DB
}

func (*Config) GetType() string {
	return Type
}

func (c *Config) GetDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName:           c.Driver,
		DSN:                  c.GetDNS(),
		PreferSimpleProtocol: c.PreferSimpleProtocol,
		Conn:                 c.Conn,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
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
	options := strings.Join(c.Options, " ")
	dsn := "host=" + c.Host + " port=" + c.Port + " user=" + c.Username + " password=" + c.Password + " dbname=" + c.Database + " " + options
	return dsn
}
