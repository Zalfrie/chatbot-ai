package models

import (
	"time"
)

type User struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	Role      string    `db:"role" json:"role"` // "user" or "admin"
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	// Fields for password reset
	ResetToken  string    `db:"reset_token" json:"-"`
	TokenExpiry time.Time `db:"token_expiry" json:"-"`
}
