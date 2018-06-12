package easyalert

import (
	"time"
)

// UserRepository wraps all CRUD operations for users
type UserRepository interface {
	FindUser(ID uint) (User, error)
	FindUsers() ([]User, error)
	CreateUser(user User) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(user User) error
}

// User defines all fields of the user model
type User struct {
	ID             uint
	Email          string
	PasswordDigest string
	Token          string
	Admin          bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
