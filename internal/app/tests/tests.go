package tests

import (
	"flag"
	"testing"

	"github.com/google/wire"
	"golang.org/x/exp/slog"
)

const appName = "go-scaffold-test"

var (
	logLevel  string // log level
	logFormat string // log output format
	dbDriver  string // database driver
)

var ProviderSet = wire.NewSet(
	NewLogger,
	NewDB,
	NewRDB,
)

func init() {
	testing.Init()

	flag.StringVar(&logLevel, "log.level", "silent", "log level（silent、debug、info、warn、error）")
	flag.StringVar(&logFormat, "log.format", "json", "log output format（text、json）")
	flag.StringVar(&dbDriver, "database.driver", "mysql", "database driver")

	flag.Parse()
}

type Tests struct {
	Logger *slog.Logger
	DB     *DB
	Redis  *Redis
}

func New(
	logger *slog.Logger,
	db *DB,
	redis *Redis,
) *Tests {
	return &Tests{
		Logger: logger,
		DB:     db,
		Redis:  redis,
	}
}
