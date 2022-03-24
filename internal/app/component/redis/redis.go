package redis

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"go-scaffold/internal/app/config"
	"time"
)

type Config struct {
	Host               string
	Port               string
	Username           string
	Password           string
	DB                 int
	MaxRetries         int
	MinRetryBackoff    time.Duration
	MaxRetryBackoff    time.Duration
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
}

func NewConfig(rdbConfig *config.Redis) *Config {
	if rdbConfig == nil {
		return nil
	}

	return &Config{
		Host:               rdbConfig.Host,
		Port:               rdbConfig.Port,
		Username:           rdbConfig.Username,
		Password:           rdbConfig.Password,
		DB:                 rdbConfig.DB,
		MaxRetries:         rdbConfig.MaxRetries,
		MinRetryBackoff:    rdbConfig.MinRetryBackoff,
		MaxRetryBackoff:    rdbConfig.MaxRetryBackoff,
		DialTimeout:        rdbConfig.DialTimeout,
		ReadTimeout:        rdbConfig.ReadTimeout,
		WriteTimeout:       rdbConfig.WriteTimeout,
		PoolSize:           rdbConfig.PoolSize,
		MinIdleConns:       rdbConfig.MinIdleConns,
		MaxConnAge:         rdbConfig.MaxConnAge,
		PoolTimeout:        rdbConfig.PoolTimeout,
		IdleTimeout:        rdbConfig.IdleTimeout,
		IdleCheckFrequency: rdbConfig.IdleCheckFrequency,
	}
}

// New 创建 redis 客户端
func New(config *Config, logger log.Logger) (*redis.Client, func(), error) {
	if config == nil {
		return nil, func() {}, nil
	}

	option := &redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.Host, config.Port),
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
		option.MinRetryBackoff = config.MinRetryBackoff * time.Second
	}
	if config.MaxRetryBackoff != 0 {
		option.MaxRetryBackoff = config.MaxRetryBackoff * time.Second
	}
	if config.DialTimeout != 0 {
		option.DialTimeout = config.DialTimeout * time.Second
	}
	if config.ReadTimeout != 0 {
		option.ReadTimeout = config.ReadTimeout * time.Second
	}
	if config.WriteTimeout != 0 {
		option.WriteTimeout = config.WriteTimeout * time.Second
	}
	if config.PoolSize != 0 {
		option.PoolSize = config.PoolSize
	}
	if config.MinIdleConns != 0 {
		option.MinIdleConns = config.MinIdleConns
	}
	if config.MaxConnAge != 0 {
		option.MaxConnAge = config.MaxConnAge * time.Second
	}
	if config.PoolTimeout != 0 {
		option.PoolTimeout = config.PoolTimeout * time.Second
	}
	if config.IdleTimeout != 0 {
		option.IdleTimeout = config.IdleTimeout * time.Second
	}
	if config.IdleCheckFrequency != 0 {
		option.IdleCheckFrequency = config.IdleCheckFrequency * time.Second
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
