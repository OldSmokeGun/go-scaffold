package redis

import (
	"context"
	"time"

	"go-scaffold/internal/config"

	"github.com/go-redis/redis/v8"
)

// New build redis client
func New(ctx context.Context, conf config.Redis) (*redis.Client, error) {
	option := &redis.Options{
		Addr: conf.Addr,
	}
	if conf.Username != "" {
		option.Username = conf.Username
	}
	if conf.Password != "" {
		option.Password = conf.Password
	}
	if conf.DB != 0 {
		option.DB = conf.DB
	}
	if conf.MaxRetries != 0 {
		option.MaxRetries = conf.MaxRetries
	}
	if conf.MinRetryBackoff != 0 {
		option.MinRetryBackoff = conf.MinRetryBackoff * time.Second
	}
	if conf.MaxRetryBackoff != 0 {
		option.MaxRetryBackoff = conf.MaxRetryBackoff * time.Second
	}
	if conf.DialTimeout != 0 {
		option.DialTimeout = conf.DialTimeout * time.Second
	}
	if conf.ReadTimeout != 0 {
		option.ReadTimeout = conf.ReadTimeout * time.Second
	}
	if conf.WriteTimeout != 0 {
		option.WriteTimeout = conf.WriteTimeout * time.Second
	}
	if conf.PoolSize != 0 {
		option.PoolSize = conf.PoolSize
	}
	if conf.MinIdleConns != 0 {
		option.MinIdleConns = conf.MinIdleConns
	}
	if conf.MaxConnAge != 0 {
		option.MaxConnAge = conf.MaxConnAge * time.Second
	}
	if conf.PoolTimeout != 0 {
		option.PoolTimeout = conf.PoolTimeout * time.Second
	}
	if conf.IdleTimeout != 0 {
		option.IdleTimeout = conf.IdleTimeout * time.Second
	}
	if conf.IdleCheckFrequency != 0 {
		option.IdleCheckFrequency = conf.IdleCheckFrequency * time.Second
	}

	client := redis.NewClient(option)

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}
