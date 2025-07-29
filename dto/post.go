package dto

import "time"

type Post struct {
	ID        int64     `db:"id"`
	Content   string    `db:"content"`
	UserID    int64     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}
