package db

import (
	"errors"
	"gin-scaffold/internal/db/mysql"
	"gin-scaffold/internal/db/postgres"
	"gorm.io/gorm"
)

var (
	ErrConversionType = errors.New("conversion type failed")
)

type Config interface {
	GetType() string
	GetDB() (*gorm.DB, error)
}

func Init(c Config) (*gorm.DB, error) {
	var (
		err error
		db  *gorm.DB
	)

	switch c.GetType() {
	case "mysql":
		config, ok := c.(*mysql.Config)
		if !ok {
			return nil, ErrConversionType
		}

		db, err = config.GetDB()
		if err != nil {
			return nil, err
		}
	case "postgres":
		config, ok := c.(*postgres.Config)
		if !ok {
			return nil, ErrConversionType
		}

		db, err = config.GetDB()
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported database type " + c.GetType())
	}

	return db, nil
}
