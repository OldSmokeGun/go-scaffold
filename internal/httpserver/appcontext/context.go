package appcontext

import (
	"gin-scaffold/internal/httpserver/appconfig"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Context 是传递给应用的依赖
type Context struct {
	config      *appconfig.Config
	logger      *logrus.Logger
	db          *gorm.DB
	redisClient *redis.Client
}

func New() *Context {
	return new(Context)
}

// SetConfig 设置日志对象
func (c *Context) SetConfig(config *appconfig.Config) {
	c.config = config
}

// Config 获取日志对象
func (c *Context) Config() *appconfig.Config {
	return c.config
}

// SetLogger 设置日志对象
func (c *Context) SetLogger(logger *logrus.Logger) {
	c.logger = logger
}

// Logger 获取日志对象
func (c *Context) Logger() *logrus.Logger {
	return c.logger
}

// SetDB 设置数据库对象
func (c *Context) SetDB(db *gorm.DB) {
	c.db = db
}

// DB 获取数据库对象
func (c *Context) DB() *gorm.DB {
	return c.db
}

// SetRedisClient 设置 redis 客户端对象
func (c *Context) SetRedisClient(rc *redis.Client) {
	c.redisClient = rc
}

// RedisClient 获取 redis 客户端对象
func (c *Context) RedisClient() *redis.Client {
	return c.redisClient
}
