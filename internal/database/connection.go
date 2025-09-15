package database

import (
	"golang-final-test/internal/config"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

var (
	db   *gorm.DB
	once sync.Once
)

func InitDB() *gorm.DB {
	once.Do(func() {
		cfg := config.LoadDB()
		var err error
		db, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}
	})
	return db
}

func WithTransaction(fn func(tx *gorm.DB) error) error {
	return InitDB().Transaction(fn)
}
