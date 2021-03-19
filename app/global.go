package app

import "github.com/go-redis/redis"

var redisClient *redis.Client

func SetRedisClient(r *redis.Client) {
	redisClient = r
}

func RedisClient() *redis.Client {
	return redisClient
}
