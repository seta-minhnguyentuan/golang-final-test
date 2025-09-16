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

func LoadDB() *DatabaseConfig {
	_ = godotenv.Load()

	port, err := strconv.Atoi(utils.GetEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	return &DatabaseConfig{
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		User:     utils.GetEnv("DB_USER", "postgres"),
		Password: utils.GetEnv("DB_PASSWORD", ""),
		DBName:   utils.GetEnv("DB_NAME", "golang_final_test"),
		Port:     port,
		SSLMode:  utils.GetEnv("DB_SSLMODE", "disable"),
	}

}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode,
	)
}
