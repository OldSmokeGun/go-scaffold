package tests

import (
	"log/slog"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

// Redis client wrapper
type Redis struct {
	client *redis.Client
	Mock   redismock.ClientMock
}

// NewRDB init redis client
func NewRDB(logger *slog.Logger) (*Redis, func(), error) {
	client, mock := redismock.NewClientMock()

	cleanup := func() {
		if err := client.Close(); err != nil {
			logger.Error("close redis client", slog.Any("error", err))
		}
	}

	return &Redis{
		client: client,
		Mock:   mock,
	}, cleanup, nil
}
