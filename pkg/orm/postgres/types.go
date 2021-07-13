package postgres

import (
	"database/sql"
	"gorm.io/gorm/logger"
	"time"
)

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
	Logger               logger.Interface
	Conn                 *sql.DB
	PreferSimpleProtocol bool
}
