package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Config struct {
	Addr     string
	Port     int
	Username string
	Password string
	Database int
}

func Initialize(config *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Addr, config.Port),
		Password: config.Password,
		DB:       config.Database,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}
