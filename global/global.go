package global

import (
	"gin-scaffold/internal/web/config"
	"github.com/go-redis/redis/v8"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Global 是传递给应用的依赖
var (
	logRotate   *rotatelogs.RotateLogs
	conf        *config.Config
	logger      *logrus.Logger
	db          *gorm.DB
	redisClient *redis.Client
)

// SetLogRotate 设置日志轮转对象
func SetLogRotate(lr *rotatelogs.RotateLogs) {
	logRotate = lr
}

// LogRotate 获取日志轮转对象
func LogRotate() *rotatelogs.RotateLogs {
	return logRotate
}

// SetConfig 设置日志对象
func SetConfig(c *config.Config) {
	conf = c
}

// Config 获取日志对象
func Config() *config.Config {
	return conf
}

// SetLogger 设置日志对象
func SetLogger(l *logrus.Logger) {
	logger = l
}

// Logger 获取日志对象
func Logger() *logrus.Logger {
	return logger
}

// SetDB 设置数据库对象
func SetDB(d *gorm.DB) {
	db = d
}

// DB 获取数据库对象
func DB() *gorm.DB {
	return db
}

// SetRedisClient 设置 redis 客户端对象
func SetRedisClient(rc *redis.Client) {
	redisClient = rc
}

// RedisClient 获取 redis 客户端对象
func RedisClient() *redis.Client {
	return redisClient
}
