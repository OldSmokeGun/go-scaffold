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
func New(ctx context.Context, conf config.Database, logger *slog.Logger) (gdb *gorm.DB, err error) {
	sdb, err := db.New(ctx, conf.DatabaseConn)
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

func registerResolver(ctx context.Context, gdb *gorm.DB, conf config.Database) error {
	rcs := conf.Resolvers
	conn := conf.DatabaseConn

	resolvers := make([]*config.DatabaseResolver, 0, len(rcs))
	for _, rc := range rcs {
		rc.Driver = conn.Driver

		if rc.MaxIdleConn == 0 {
			rc.MaxIdleConn = conn.MaxIdleConn
		}
		if rc.MaxOpenConn == 0 {
			rc.MaxOpenConn = conn.MaxOpenConn
		}
		if rc.ConnMaxIdleTime == 0 {
			rc.ConnMaxIdleTime = conn.ConnMaxIdleTime
		}
		if rc.ConnMaxLifeTime == 0 {
			rc.ConnMaxLifeTime = conn.ConnMaxLifeTime
		}

		resolvers = append(resolvers, rc)
	}

	plugin, err := buildResolver(ctx, resolvers)
	if err != nil {
		return err
	}
	return gdb.Use(plugin)
}

func buildResolver(ctx context.Context, resolvers []*config.DatabaseResolver) (gorm.Plugin, error) {
	var (
		sources  = make([]gorm.Dialector, 0, len(resolvers))
		replicas = make([]gorm.Dialector, 0, len(resolvers))
	)

	for _, resolver := range resolvers {
		sdb, err := db.New(ctx, resolver.DatabaseConn)
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
