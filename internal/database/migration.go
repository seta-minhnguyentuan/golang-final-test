package database

import (
	"golang-final-test/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Post{}, &models.ActivityLog{}); err != nil {
		return err
	}

	if err := db.Exec(`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_posts_tags_gin ON posts USING gin(tags);`).Error; err != nil {
		return err
	}

	return nil
}
