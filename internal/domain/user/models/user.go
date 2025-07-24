package models

import "time"

type User struct {
	ID        uint
	Password  string // Hashed!
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserCredentials struct {
	Email    string
	Password string
}
