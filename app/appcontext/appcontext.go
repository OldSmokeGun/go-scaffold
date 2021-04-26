package appcontext

import (
	"gin-scaffold/app/appconfig"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Context 是 app 的上下文依赖
type Context struct {
	Config       appconfig.Config
	DB           *gorm.DB
	Logger       *logrus.Logger
	Configurator *viper.Viper
	RedisClient  *redis.Client
}

// OptionFunc 选项模式中的回调方法
type OptionFunc func(ctx *Context)

// WithConfig 传递 config.Config 对象到 Context 中
func WithConfig(config appconfig.Config) OptionFunc {
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

// WithLogger 传递 *logrus.Logger 对象到 Context 中
func WithLogger(logger *logrus.Logger) OptionFunc {
	return func(ctx *Context) {
		ctx.Logger = logger
	}
}

// WithConfigurator 传递 *viper.Viper 对象到 Context 中
func WithConfigurator(configurator *viper.Viper) OptionFunc {
	return func(ctx *Context) {
		ctx.Configurator = configurator
	}
}

// WithRedisClient 传递 *redis.Client 对象到 Context 中
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
