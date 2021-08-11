package appcontext

import (
	"gin-scaffold/internal/httpserver/appconfig"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Context 是传递给应用的依赖
type Context struct {
	logRotate   *rotatelogs.RotateLogs
	config      *appconfig.Config
	logger      *logrus.Logger
	db          *gorm.DB
	redisClient *redis.Client
	router      *gin.Engine
}

func New() *Context {
	return new(Context)
}

// SetLogRotate 设置日志轮转对象
func (c *Context) SetLogRotate(logRotate *rotatelogs.RotateLogs) {
	c.logRotate = logRotate
}

// LogRotate 获取日志轮转对象
func (c *Context) LogRotate() *rotatelogs.RotateLogs {
	return c.logRotate
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

// SetRouter 设置路由对象
func (c *Context) SetRouter(router *gin.Engine) {
	c.router = router
}

// Router 获取路由对象
func (c *Context) Router() *gin.Engine {
	return c.router
}
