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
)

type User struct {
	UserId    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`

	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

func (u *User) Update(uu *User) {
	if uu.Username != "" {
		u.Username = uu.Username
	}
	if uu.Name != "" {
		u.Name = uu.Name
	}
	if uu.Email != "" {
		u.Email = uu.Email
	}
}

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
}
