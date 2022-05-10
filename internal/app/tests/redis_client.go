package tests

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

// RedisClient 客户端
type RedisClient struct {
	DB   *redis.Client
	Mock redismock.ClientMock
}

// NewRDB 初始化测试 redis 客户端
func NewRDB(kLogger klog.Logger) (*RedisClient, func(), error) {
	logger := klog.NewHelper(kLogger)

	rdb, mock := redismock.NewClientMock()

	cleanup := func() {
		if err := rdb.Close(); err != nil {
			logger.Error(err)
		}
	}

	return &RedisClient{
		DB:   rdb,
		Mock: mock,
	}, cleanup, nil
}
