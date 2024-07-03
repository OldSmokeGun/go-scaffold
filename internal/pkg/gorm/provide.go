package gorm

import (
	"context"
	"log/slog"

	"gorm.io/gorm"

	"go-scaffold/internal/config"
)

type DefaultDB = gorm.DB

// ProvideDefault default gorm
func ProvideDefault(ctx context.Context, conf config.DefaultDatabase, logger *slog.Logger) (db *DefaultDB, cleanup func(), err error) {
	db, err = New(ctx, conf, logger)
	if err != nil {
		return nil, nil, err
	}

	cleanup = func() {
		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}

		if err := sqlDB.Close(); err != nil {
			panic(err)
		}
	}

	return db, cleanup, nil
}
