package ent

import (
	"context"
	"go-scaffold/internal/app/component/ent/ent"
	"go-scaffold/internal/app/component/ent/ent/migrate"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
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
	Addr            string
	Database        string
	Username        string
	Password        string
	Options         string
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxIdleTime int64
	ConnMaxLifeTime int64
}

func New(config *Config, logger log.Logger) (*ent.Client, func(), error) {
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

	hLogger := log.NewHelper(logger)
	cleanup := func() {
		hLogger.Info("closing the ent resources")

		err = driver.Close()
		if err != nil {
			hLogger.Error(err)
		}
	}

	client := ent.NewClient(
		ent.Driver(driver),
		ent.Log(func(i ...any) {
			hLogger.Debug(i)
		}),
	)

	if err := client.Schema.Create(
		context.Background(),
		migrate.WithForeignKeys(false),
	); err != nil {
		hLogger.Errorf("failed to creat schema resources: %v", err)
		return nil, nil, err
	}

	return client, cleanup, nil
}

func buildSource(c *Config) string {
	dsn := ""

	switch c.Driver {
	case PostgresSQL:
		var host, port string
		s := strings.SplitN(c.Addr, ":", 2)
		if len(s) == 2 {
			host = s[0]
			port = s[1]
		}

		dsn = "host=" + host + " port=" + port + " user=" + c.Username + " password=" + c.Password + " dbname=" + c.Database + " " + c.Options
	case MySQL:
		fallthrough
	default:
		dsn = c.Username + ":" + c.Password + "@tcp(" + c.Addr + ")/" + c.Database + "?" + c.Options
	}

	return dsn
}
