package models

import (
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Action   string    `gorm:"column:action" json:"action"`
	PostID   uuid.UUID `gorm:"type:uuid" json:"post_id"`
	LoggedAt time.Time `gorm:"autoCreateTime" json:"logged_at"`
}
