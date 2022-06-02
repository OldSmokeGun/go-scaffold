package tests

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"go.uber.org/zap"
)

// RedisClient 客户端
type RedisClient struct {
	DB   *redis.Client
	Mock redismock.ClientMock
}

// NewRDB 初始化测试 redis 客户端
func NewRDB(logger *zap.Logger) (*RedisClient, func(), error) {
	rdb, mock := redismock.NewClientMock()

	cleanup := func() {
		if err := rdb.Close(); err != nil {
			logger.Sugar().Error(err)
		}
	}

	return &RedisClient{
		DB:   rdb,
		Mock: mock,
	}, cleanup, nil
}
