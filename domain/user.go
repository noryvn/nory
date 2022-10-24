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
	UserId    string    `json:"userId"`    // immutable, unique
	CreatedAt time.Time `json:"createdAt"` // immutable

	Username string `json:"username"` //  mutable, unique
	Name     string `json:"name"`     // mutable
	Email    string `json:"email"`    // immutable, unique
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
