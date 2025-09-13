package config

import (
	"fmt"
	"golang-final-test/pkg/utils"
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var cfg *DatabaseConfig

func LoadDB() *DatabaseConfig {
	_ = godotenv.Load()

	if cfg != nil {
		return cfg
	}

	port, err := strconv.Atoi(utils.GetEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	cfg = &DatabaseConfig{
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		User:     utils.GetEnv("DB_USER", "postgres"),
		Password: utils.GetEnv("DB_PASSWORD", ""),
		DBName:   utils.GetEnv("DB_NAME", "postgres"),
		Port:     port,
		SSLMode:  utils.GetEnv("DB_SSLMODE", "disable"),
	}

	return cfg
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode,
	)
}
