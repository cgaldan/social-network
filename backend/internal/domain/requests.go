package domain

import "time"

type RegisterRequest struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Nickname    string    `json:"nickname"`
	Gender      string    `json:"gender"`
	AvatarPath  string    `json:"avatar_path"`
	AboutMe     string    `json:"about_me"`
	IsPublic    bool      `json:"is_public"`
}

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}
