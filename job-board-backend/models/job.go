package models

import "database/sql"

type Job struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	UserID      sql.NullInt64 `json:"user_id"`
}
