package data

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"go-scaffold/internal/app/component/data/ent"
	"go-scaffold/internal/app/component/data/ent/migrate"
	"go-scaffold/internal/app/config"
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
	Port            string
	Database        string
	Username        string
	Password        string
	Options         []string
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxIdleTime int64
	ConnMaxLifeTime int64
}

func NewConfig(dbConfig *config.DB) *Config {
	if dbConfig == nil {
		return nil
	}

	return &Config{
		Driver:          Driver(dbConfig.Driver),
		Host:            dbConfig.Host,
		Port:            dbConfig.Port,
		Database:        dbConfig.Database,
		Username:        dbConfig.Username,
		Password:        dbConfig.Password,
		Options:         dbConfig.Options,
		MaxIdleConn:     dbConfig.MaxIdleConn,
		MaxOpenConn:     dbConfig.MaxOpenConn,
		ConnMaxIdleTime: dbConfig.ConnMaxIdleTime,
		ConnMaxLifeTime: dbConfig.ConnMaxLifeTime,
	}
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
		ent.Log(func(i ...interface{}) {
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
		options := strings.Join(c.Options, " ")
		dsn = "host=" + c.Host + " port=" + c.Port + " user=" + c.Username + " password=" + c.Password + " dbname=" + c.Database + " " + options
	case MySQL:
		fallthrough
	default:
		options := strings.Join(c.Options, "&")
		dsn = c.Username + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.Database + "?" + options
	}

	return dsn
}
