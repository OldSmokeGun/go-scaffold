package config

import (
	"net/url"
	"time"

	"github.com/samber/lo"
)

var supportedDrivers = []DBDriver{MySQL, Postgres}

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
	DBDsn
	DBConnPool
}

// DBDsn database dsn config
type DBDsn struct {
	Driver   DBDriver `json:"driver"`
	Addr     string   `json:"addr"`
	Database string   `json:"database"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Options  string   `json:"options"`
}

// EnableMultiStatement enable the execution of multi sql statement
func (d *DBDsn) EnableMultiStatement() error {
	options, err := url.ParseQuery(d.Options)
	if err != nil {
		return err
	}

	if !options.Has("multiStatements") {
		options.Set("multiStatements", "true")
	}
	
	d.Options = options.Encode()
	return nil
}

// DBConnPool database connection pool config
type DBConnPool struct {
	MaxIdleConn     int           `json:"maxIdleConn"`
	MaxOpenConn     int           `json:"maxOpenConn"`
	ConnMaxIdleTime time.Duration `json:"connMaxIdleTime"`
	ConnMaxLifeTime time.Duration `json:"connMaxLifeTime"`
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
