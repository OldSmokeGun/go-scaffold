package gorm

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"go-scaffold/internal/app/pkg/db"
	"go-scaffold/internal/config"
	glog "go-scaffold/pkg/log/gorm"
)

// ErrUnsupportedResolverType unsupported resolver type
var ErrUnsupportedResolverType = errors.New("unsupported resolver type")

// New build gorm
func New(ctx context.Context, conf config.DB, logger *slog.Logger) (gdb *gorm.DB, err error) {
	sdb, err := db.New(ctx, conf.DBConn)
	if err != nil {
		return nil, err
	}

	if logger != nil {
		gormlogger.Default = glog.NewLogger(logger, glog.Config{
			SlowThreshold:             200 * time.Millisecond,
			IgnoreRecordNotFoundError: false,
			LogInfo:                   conf.LogInfo,
		})
	}

	dialect, err := NewDialect(conf.Driver, sdb)
	if err != nil {
		return nil, err
	}

	gdb, err = gorm.Open(dialect, &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err = registerResolver(ctx, gdb, conf); err != nil {
		return nil, err
	}

	return gdb, nil
}

func registerResolver(ctx context.Context, gdb *gorm.DB, conf config.DB) error {
	rcs := conf.Resolvers
	conn := conf.DBConn

	resolvers := make([]*config.DBResolver, 0, len(rcs))
	for _, rc := range rcs {
		rc.Driver = conn.Driver

		rcp := conn.DBConnPool
		if rc.MaxIdleConn > 0 {
			rcp.MaxIdleConn = rc.MaxIdleConn
		}
		if rc.MaxOpenConn > 0 {
			rcp.MaxOpenConn = rc.MaxOpenConn
		}
		if rc.ConnMaxIdleTime > 0 {
			rcp.ConnMaxIdleTime = rc.ConnMaxIdleTime
		}
		if rc.ConnMaxLifeTime > 0 {
			rcp.ConnMaxLifeTime = rc.ConnMaxLifeTime
		}

		resolvers = append(resolvers, &config.DBResolver{
			Type: rc.Type,
			DBConn: config.DBConn{
				DBDsn:      rc.DBDsn,
				DBConnPool: rcp,
			},
		})
	}

	plugin, err := buildResolver(ctx, resolvers)
	if err != nil {
		return err
	}
	return gdb.Use(plugin)
}

func buildResolver(ctx context.Context, resolvers []*config.DBResolver) (gorm.Plugin, error) {
	var (
		sources  = make([]gorm.Dialector, 0, len(resolvers))
		replicas = make([]gorm.Dialector, 0, len(resolvers))
	)

	for _, resolver := range resolvers {
		sdb, err := db.New(ctx, resolver.DBConn)
		if err != nil {
			return nil, err
		}

		dialect, err := NewDialect(resolver.Driver, sdb)
		if err != nil {
			return nil, err
		}

		switch resolver.Type {
		case config.Source:
			sources = append(sources, dialect)
		case config.Replica:
			replicas = append(replicas, dialect)
		default:
			return nil, ErrUnsupportedResolverType
		}
	}

	return dbresolver.Register(dbresolver.Config{
		Sources:  sources,
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	}), nil
}
