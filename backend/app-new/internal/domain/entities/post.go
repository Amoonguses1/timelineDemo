package entities

import (
	"time"

	"github.com/google/uuid"
)

// Post represents an entry of `posts` table.
type Post struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
