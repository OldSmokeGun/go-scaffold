package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

type Config struct {
	Host               string
	Port               int
	Username           string
	Password           string
	DB                 int
	MaxRetries         int
	MinRetryBackoff    int64
	MaxRetryBackoff    int64
	DialTimeout        int64
	ReadTimeout        int64
	WriteTimeout       int64
	PoolSize           int
	MinIdleConns       int
	MaxConnAge         int64
	PoolTimeout        int64
	IdleTimeout        int64
	IdleCheckFrequency int64
}

// New 创建 redis 客户端
func New(config *Config, logger log.Logger) (*redis.Client, func(), error) {
	if config == nil {
		return nil, func() {}, nil
	}

	option := &redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
	}
	if config.Username != "" {
		option.Username = config.Username
	}
	if config.Password != "" {
		option.Password = config.Password
	}
	if config.DB != 0 {
		option.DB = config.DB
	}
	if config.MaxRetries != 0 {
		option.MaxRetries = config.MaxRetries
	}
	if config.MinRetryBackoff != 0 {
		option.MinRetryBackoff = time.Duration(config.MinRetryBackoff) * time.Second
	}
	if config.MaxRetryBackoff != 0 {
		option.MaxRetryBackoff = time.Duration(config.MaxRetryBackoff) * time.Second
	}
	if config.DialTimeout != 0 {
		option.DialTimeout = time.Duration(config.DialTimeout) * time.Second
	}
	if config.ReadTimeout != 0 {
		option.ReadTimeout = time.Duration(config.ReadTimeout) * time.Second
	}
	if config.WriteTimeout != 0 {
		option.WriteTimeout = time.Duration(config.WriteTimeout) * time.Second
	}
	if config.PoolSize != 0 {
		option.PoolSize = config.PoolSize
	}
	if config.MinIdleConns != 0 {
		option.MinIdleConns = config.MinIdleConns
	}
	if config.MaxConnAge != 0 {
		option.MaxConnAge = time.Duration(config.MaxConnAge) * time.Second
	}
	if config.PoolTimeout != 0 {
		option.PoolTimeout = time.Duration(config.PoolTimeout) * time.Second
	}
	if config.IdleTimeout != 0 {
		option.IdleTimeout = time.Duration(config.IdleTimeout) * time.Second
	}
	if config.IdleCheckFrequency != 0 {
		option.IdleCheckFrequency = time.Duration(config.IdleCheckFrequency) * time.Second
	}

	client := redis.NewClient(option)

	ctx := context.Background()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		log.NewHelper(logger).Info("closing the redis client")

		if err := client.Close(); err != nil {
			log.NewHelper(logger).Error(err.Error())
		}
	}

	return client, cleanup, nil
}
