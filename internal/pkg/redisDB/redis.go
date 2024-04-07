package redisDB

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// Define a Redis client
var RedisClient *redis.Client

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis server address
		Password: "",           // No password
		DB:       0,            // Use default DB
	})
	// Ping the Redis server to check connection
	ctx := context.Background()
	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)
}

func GetClient() (*redis.Client, error) {
	if RedisClient == nil {
		return nil, fmt.Errorf("redis client is not initialized")
	}
	return RedisClient, nil
}

func GetCacheData(ctx context.Context, key string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("redis client is not initialized")
	}
	return RedisClient.Get(ctx, key).Result()
}

func SetCacheData(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	return RedisClient.Set(ctx, key, value, expiration).Err()
}
