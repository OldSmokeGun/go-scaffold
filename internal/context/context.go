package context

import (
	"gin-scaffold/internal/config"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Context 是传递给应用的依赖
type Context struct {
	Config config.Config

	rootCommand  *cobra.Command
	logger       *logrus.Logger
	configurator *viper.Viper
	db           *gorm.DB
	redisClient  *redis.Client
}

var ctx = new(Context)

func Default() *Context {
	return ctx
}

// SetRootCommand 设置根命令
func (*Context) SetRootCommand(cmd *cobra.Command) {
	ctx.rootCommand = cmd
}

// GetRootCommand 获取根命令
func (*Context) GetRootCommand() *cobra.Command {
	return ctx.rootCommand
}

// SetConfigurator 设置配置对象
func (*Context) SetConfigurator(cfg *viper.Viper) {
	ctx.configurator = cfg
}

// GetConfigurator 获取配置对象
func (*Context) GetConfigurator() *viper.Viper {
	return ctx.configurator
}

// SetLogger 设置日志对象
func (*Context) SetLogger(logger *logrus.Logger) {
	ctx.logger = logger
}

// GetLogger 获取日志对象
func (*Context) GetLogger() *logrus.Logger {
	return ctx.logger
}

// SetDB 设置数据库对象
func (*Context) SetDB(db *gorm.DB) {
	ctx.db = db
}

// GetDB 获取数据库对象
func (*Context) GetDB() *gorm.DB {
	return ctx.db
}

// SetRedisClient 设置 redis 客户端对象
func (*Context) SetRedisClient(rc *redis.Client) {
	ctx.redisClient = rc
}

// GetRedisClient 获取 redis 客户端对象
func (*Context) GetRedisClient() *redis.Client {
	return ctx.redisClient
}
