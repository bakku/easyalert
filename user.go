package easyalert

import (
	"time"

	"golang.org/x/crypto/bcrypt"
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

func (u *User) HashPassword(pass string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordDigest = string(hash)

	return nil
}

func (u *User) ValidPassword(pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(pass)) == nil
}
