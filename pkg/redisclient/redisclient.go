package redisclient

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Setup 返回 *redis.Client
func Setup(config Config) (*redis.Client, error) {
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
		option.MinRetryBackoff = config.MinRetryBackoff
	}
	if config.MaxRetryBackoff != 0 {
		option.MaxRetryBackoff = config.MaxRetryBackoff
	}
	if config.DialTimeout != 0 {
		option.DialTimeout = config.DialTimeout
	}
	if config.ReadTimeout != 0 {
		option.ReadTimeout = config.ReadTimeout
	}
	if config.WriteTimeout != 0 {
		option.WriteTimeout = config.WriteTimeout
	}
	if config.PoolSize != 0 {
		option.PoolSize = config.PoolSize
	}
	if config.MinIdleConns != 0 {
		option.MinIdleConns = config.MinIdleConns
	}
	if config.MaxConnAge != 0 {
		option.MaxConnAge = config.MaxConnAge
	}
	if config.PoolTimeout != 0 {
		option.PoolTimeout = config.PoolTimeout
	}
	if config.IdleTimeout != 0 {
		option.IdleTimeout = config.IdleTimeout
	}
	if config.IdleCheckFrequency != 0 {
		option.IdleCheckFrequency = config.IdleCheckFrequency
	}

	client := redis.NewClient(option)

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}

// MustSetup 返回 *redis.Client
func MustSetup(config Config) *redis.Client {
	client, err := Setup(config)
	if err != nil {
		panic(err)
	}

	return client
}
