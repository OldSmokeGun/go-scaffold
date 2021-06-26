package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	ConnMaxLifeTime           time.Duration
	LogLevel                  logger.Interface
	Conn                      gorm.ConnPool
	SkipInitializeWithVersion bool
	DefaultStringSize         uint
	DisableDatetimePrecision  bool
	DontSupportRenameIndex    bool
	DontSupportRenameColumn   bool
}
