package config

import (
	"net/url"
	"time"

	"github.com/samber/lo"
)

var supportedDrivers = []DatabaseDriver{MySQL, Postgres, SQLite}

// DatabaseGroup config
type DatabaseGroup struct {
	Default *DefaultDatabase `json:"default"`
}

func (DatabaseGroup) GetName() string {
	return "database"
}

// DefaultDatabase default database config
type DefaultDatabase = Database

func (DefaultDatabase) GetName() string {
	return "database.default"
}

// Database config
type Database struct {
	DatabaseConn
	LogInfo   bool                `json:"logInfo"`
	Resolvers []*DatabaseResolver `json:"resolvers"`
}

// DatabaseConn connection config
type DatabaseConn struct {
	Driver          DatabaseDriver `json:"driver"`
	DSN             string         `json:"dsn"`
	MaxIdleConn     int            `json:"maxIdleConn"`
	MaxOpenConn     int            `json:"maxOpenConn"`
	ConnMaxIdleTime time.Duration  `json:"connMaxIdleTime"`
	ConnMaxLifeTime time.Duration  `json:"connMaxLifeTime"`
}

// EnableMultiStatement enable the execution of multi sql statement
func (d *DatabaseConn) EnableMultiStatement() error {
	if d.Driver != MySQL {
		return nil
	}

	options, err := url.Parse(d.DSN)
	if err != nil {
		return err
	}

	if !options.Query().Has("multiStatements") {
		q := options.Query()
		q.Set("multiStatements", "true")
		options.RawQuery = q.Encode()
	}

	d.DSN = options.String()
	return nil
}

// DatabaseDriver database driver type
type DatabaseDriver string

func (d DatabaseDriver) String() string {
	return string(d)
}

// IsSupported check that the driver is supported
func (d DatabaseDriver) IsSupported() bool {
	return lo.Contains(supportedDrivers, d)
}

const (
	MySQL    DatabaseDriver = "mysql"
	Postgres DatabaseDriver = "pgx"
	SQLite   DatabaseDriver = "sqlite3"
)

// DatabaseResolver database resolver config
type DatabaseResolver struct {
	Type DatabaseResolverType `json:"type"`
	DatabaseConn
}

// DatabaseResolverType database resolver type
type DatabaseResolverType string

const (
	Source  DatabaseResolverType = "source"
	Replica DatabaseResolverType = "replica"
)
