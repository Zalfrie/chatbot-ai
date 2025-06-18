package models

import (
	"time"
)

type Message struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Content   string    `db:"content" json:"content"`
	Private   bool      `db:"is_private" json:"is_private"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
