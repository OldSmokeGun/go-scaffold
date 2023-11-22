package config

import (
	"net/url"
	"time"

	"github.com/samber/lo"
)

var supportedDrivers = []DBDriver{MySQL, Postgres, SQLite}

// DB database config
type DB struct {
	DBConn
	LogInfo   bool          `json:"logInfo"`
	Resolvers []*DBResolver `json:"resolvers"`
}

func (DB) GetName() string {
	return "db"
}

// DBConn database connection config
type DBConn struct {
	Driver          DBDriver      `json:"driver"`
	DSN             string        `json:"dsn"`
	MaxIdleConn     int           `json:"maxIdleConn"`
	MaxOpenConn     int           `json:"maxOpenConn"`
	ConnMaxIdleTime time.Duration `json:"connMaxIdleTime"`
	ConnMaxLifeTime time.Duration `json:"connMaxLifeTime"`
}

// EnableMultiStatement enable the execution of multi sql statement
func (d *DBConn) EnableMultiStatement() error {
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

// DBDriver database driver type
type DBDriver string

func (d DBDriver) String() string {
	return string(d)
}

// IsSupported check that the driver is supported
func (d DBDriver) IsSupported() bool {
	return lo.Contains(supportedDrivers, d)
}

const (
	MySQL    DBDriver = "mysql"
	Postgres DBDriver = "postgres"
	SQLite   DBDriver = "sqlite3"
)

// DBResolver database resolver config
type DBResolver struct {
	Type DBResolverType `json:"type"`
	DBConn
}

// DBResolverType database resolver type
type DBResolverType string

const (
	Source  DBResolverType = "source"
	Replica DBResolverType = "replica"
)
