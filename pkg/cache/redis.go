package cache

import (
	"context"
	"golang-final-test/internal/config"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

func InitRedis() *redis.Client {
	once.Do(func() {
		cfg := config.LoadRedisConfig()

		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Host + ":" + cfg.Port,
			Password: cfg.Password,
			DB:       cfg.DB,
		})

		if err := client.Ping(context.Background()).Err(); err != nil {
			log.Fatalf("failed to connect redis: %v", err)
		}
	})

	return client
}

func Set(key string, value string, ttl time.Duration) error {
	return InitRedis().Set(context.Background(), key, value, ttl).Err()
}

func Get(key string) (string, error) {
	return InitRedis().Get(context.Background(), key).Result()
}

func Delete(key string) error {
	return InitRedis().Del(context.Background(), key).Err()
}
