package tests

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"golang.org/x/exp/slog"
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
			logger.Error("close redis client", err)
		}
	}

	return &Redis{
		client: client,
		Mock:   mock,
	}, cleanup, nil
}
