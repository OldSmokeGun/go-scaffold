package ent

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"go-scaffold/internal/app/component/ent/ent"
	"go-scaffold/internal/app/component/ent/ent/migrate"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type Driver string

func (d Driver) String() string {
	return string(d)
}

const (
	MySQL       Driver = "mysql"
	PostgresSQL Driver = "postgres"
)

type Config struct {
	Driver          Driver
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	Options         []string
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxIdleTime int64
	ConnMaxLifeTime int64
}

func New(config *Config, logger *zap.Logger) (*ent.Client, func(), error) {
	driver, err := sql.Open(config.Driver.String(), buildSource(config))
	if err != nil {
		return nil, nil, err
	}

	if config.MaxIdleConn > 0 {
		driver.DB().SetMaxIdleConns(config.MaxIdleConn)
	}
	if config.MaxOpenConn > 0 {
		driver.DB().SetMaxOpenConns(config.MaxOpenConn)
	}
	if config.ConnMaxIdleTime > 0 {
		driver.DB().SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second)
	}
	if config.ConnMaxLifeTime > 0 {
		driver.DB().SetConnMaxLifetime(time.Duration(config.ConnMaxLifeTime) * time.Second)
	}

	cleanup := func() {
		logger.Info("closing the ent resources")

		err = driver.Close()
		if err != nil {
			logger.Sugar().Error(err)
		}
	}

	client := ent.NewClient(
		ent.Driver(driver),
		ent.Log(func(i ...interface{}) {
			logger.Sugar().Debug(i)
		}),
	)

	if err := client.Schema.Create(
		context.Background(),
		migrate.WithForeignKeys(false),
	); err != nil {
		logger.Sugar().Errorf("failed to creat schema resources: %v", err)
		return nil, nil, err
	}

	return client, cleanup, nil
}

func buildSource(c *Config) string {
	dsn := ""

	switch c.Driver {
	case PostgresSQL:
		options := strings.Join(c.Options, " ")
		dsn = "host=" + c.Host + " port=" + strconv.Itoa(c.Port) + " user=" + c.Username + " password=" + c.Password + " dbname=" + c.Database + " " + options
	case MySQL:
		fallthrough
	default:
		options := strings.Join(c.Options, "&")
		dsn = c.Username + ":" + c.Password + "@tcp(" + c.Host + ":" + strconv.Itoa(c.Port) + ")/" + c.Database + "?" + options
	}

	return dsn
}
