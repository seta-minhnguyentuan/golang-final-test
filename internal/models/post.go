package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Post struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Title     string         `gorm:"column:title" json:"title"`
	Content   string         `gorm:"column:content" json:"content"`
	Tags      pq.StringArray `gorm:"column:tags;type:text[]" json:"tags"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}
