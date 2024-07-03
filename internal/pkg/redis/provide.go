package redis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"go-scaffold/internal/config"
)

type DefaultRedis = redis.Client

// ProvideDefault default redis client
func ProvideDefault(ctx context.Context, conf config.Redis) (*DefaultRedis, func(), error) {
	client, err := New(ctx, conf)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := client.Close(); err != nil {
			panic(err)
		}
	}

	return client, cleanup, nil
}
