package domain

import (
	"context"
	"errors"
	"time"
)

var (
	// user not found
	ErrUserNotFound = errors.New("can`t find user")
	// used for duplicate UserId
	ErrUserExists = errors.New("user already exists")
	// used for duplicate data, username, or unique sql column etc
	ErrDuplicateUser = errors.New("duplicate user data")
)

type User struct {
	UserId    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`

	Username string `json:"username"`
	Name     string `json:"name"`
	Email string `json:"email"`
}

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
}
