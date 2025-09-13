package models

import (
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	ID       uuid.UUID `json:"id"`
	Action   string    `json:"action"`
	PostID   uuid.UUID `json:"post_id"`
	LoggedAt time.Time `json:"logged_at"`
}
