package domain

import (
	"context"
	"errors"
	"time"
)

var (
	// user not found
	ErrUserNotFound  = errors.New("can`t find user")
	// duplicate UserId
	ErrUserExists    = errors.New("user already exists")
	// duplicate data, for unique sql column etc
	ErrDuplicateUser = errors.New("duplicate user data")
)

type User struct {
	UserId    string
	CreatedAt time.Time

	Username string
	Name     string
}

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
}
