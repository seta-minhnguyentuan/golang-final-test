package config

import "golang-final-test/pkg/utils"

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func LoadRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     utils.GetEnv("REDIS_HOST", "localhost"),
		Port:     utils.GetEnv("REDIS_PORT", "6379"),
		Password: utils.GetEnv("REDIS_PASSWORD", ""),
		DB:       0,
	}
}
