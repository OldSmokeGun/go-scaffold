package tests

import (
	"flag"
	"os"
	"testing"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var (
	appName     = "go-scaffold-test"
	hostname, _ = os.Hostname()
)

var (
	logLevel      string // 日志等级
	logFormat     string // 日志输出格式
	logCallerSkip int    // 日志 caller 跳过层数
)

var ProviderSet = wire.NewSet(
	NewZapLogger,
	NewLogger,
	NewDB,
	NewRDB,
)

func init() {
	testing.Init()

	flag.StringVar(&logLevel, "log.level", "silent", "日志等级（silent, debug、info、warn、error、panic、fatal）")
	flag.StringVar(&logFormat, "log.format", "json", "日志输出格式（text、json）")
	flag.IntVar(&logCallerSkip, "log.caller-skip", 4, "日志 caller 跳过层数")

	flag.Parse()
}

type Tests struct {
	ZapLogger   *zap.Logger
	Logger      klog.Logger
	DB          *DB
	RedisClient *RedisClient
}

func New(
	zapLogger *zap.Logger,
	logger klog.Logger,
	db *DB,
	redisClient *RedisClient,
) *Tests {
	return &Tests{
		ZapLogger:   zapLogger,
		Logger:      logger,
		DB:          db,
		RedisClient: redisClient,
	}
}
