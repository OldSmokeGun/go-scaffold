package gorm

import (
	"context"
	"log/slog"

	"gorm.io/gorm"

	"go-scaffold/internal/config"
)

// Provide gorm
func Provide(ctx context.Context, conf config.Database, logger *slog.Logger) (db *gorm.DB, cleanup func(), err error) {
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
