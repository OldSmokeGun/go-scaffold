package redis

import (
	"context"

	"go-scaffold/internal/config"

	"github.com/go-redis/redis/v8"
)

// Provide redis client
func Provide(ctx context.Context, conf config.Redis) (*redis.Client, func(), error) {
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
