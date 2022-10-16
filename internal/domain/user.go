package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUserNotFound  = errors.New("can`t find user")
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
