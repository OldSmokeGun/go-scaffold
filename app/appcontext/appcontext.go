package appcontext

import (
	"gin-scaffold/app/config"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// Context 是 app 的上下文依赖
type Context struct {
	Config      config.Config
	DB          *gorm.DB
	RedisClient *redis.Client
}

// OptionFunc 选项模式中的回调方法
type OptionFunc func(ctx *Context)

// WithConfig 传递 config.Config 对象到 Context 中
func WithConfig(config config.Config) OptionFunc {
	return func(ctx *Context) {
		ctx.Config = config
	}
}

// WithDB 传递 *gorm.DB 对象到 Context 中
func WithDB(db *gorm.DB) OptionFunc {
	return func(ctx *Context) {
		ctx.DB = db
	}
}

// WithRedisClient 传递 redis.Client 对象到 Context 中
func WithRedisClient(redisClient *redis.Client) OptionFunc {
	return func(ctx *Context) {
		ctx.RedisClient = redisClient
	}
}

// NewContext 返回一个 Context 对象
func NewContext(options ...OptionFunc) *Context {
	ctx := new(Context)
	for _, option := range options {
		option(ctx)
	}
	return ctx
}
