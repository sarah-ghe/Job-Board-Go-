package models

type User struct {
	ID       int
	Username string
	Email    string
	Password string // will store the hashed password
}