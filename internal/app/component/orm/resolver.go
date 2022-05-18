package orm

import (
	mysqlx "go-scaffold/internal/app/component/orm/mysql"
	postgresx "go-scaffold/internal/app/component/orm/postgres"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ResolverType resolver type
type ResolverType string

const (
	Source  ResolverType = "source"
	Replica ResolverType = "replica"
)

// Resolver provide multiple database support
type Resolver struct {
	Type ResolverType
	DSN
}

// BuildDialector build gorm.Dialector
func BuildDialector(driver Driver, dsn DSN) (dialector gorm.Dialector, err error) {
	switch driver {
	case MySQL:
		dialector = mysql.New(mysql.Config{
			DriverName: driver.String(),
			DSN: mysqlx.BuildDSN(mysqlx.Config{
				Host:     dsn.Host,
				Port:     dsn.Port,
				Database: dsn.Database,
				Username: dsn.Username,
				Password: dsn.Password,
				Options:  dsn.Options,
			}),
		})
	case Postgres:
		dialector = postgres.New(postgres.Config{
			DriverName: driver.String(),
			DSN: postgresx.BuildDSN(postgresx.Config{
				Host:     dsn.Host,
				Port:     dsn.Port,
				Database: dsn.Database,
				Username: dsn.Username,
				Password: dsn.Password,
				Options:  dsn.Options,
			}),
		})
	default:
		return nil, ErrUnsupportedType
	}

	return dialector, nil
}
