package models

import "database/sql"

type Job struct {
	ID          int
	Title       string
	Description string
	UserID      sql.NullInt64
}
