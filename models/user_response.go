package models

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}


type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}