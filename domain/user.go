package domain

import (
	"context"
	"errors"
	"time"
)

var (
	// user not found
	ErrUserNotExists = errors.New("user not found")
	// used for duplicate UserId
	ErrUserAlreadyExists = errors.New("user already exists")
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
}

type UserRepository interface {
	// create user takes an (*User) and use the UserId as id, it becaues the id came from third party authentication service
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	// update user, consider using (*User).Update(*User) to avoid overwrite immutable fields
	UpdateUser(ctx context.Context, user *User) error
}
