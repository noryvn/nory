package domain

import (
	"context"
	"errors"
	"time"
)

var (
	// user not found
	ErrUserNotExists = errors.New("user does not exists")
	// used for duplicate UserId
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserStatistics struct {
	JoinedClass int `json:"joinedClass"`
	OwnedClass  int `json:"ownedClass"`
}

type User struct {
	UserId    string    `json:"userId"`    // immutable, unique
	CreatedAt time.Time `json:"createdAt"` // immutable

	Username string `json:"username" validate:"username"`       // mutable, unique
	Name     string `json:"name" validate:"max=32"`             // mutable
	Email    string `json:"email,omitempty" validate:"max=254"` // immutable, unique

	UserStatistics *UserStatistics `json:"userStatistics,omitempty"`

	OwnedClass []*Class `json:"ownedClass,omitempty"`
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
	GetUserByUserId(ctx context.Context, id string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	// update user, consider using (*User).Update(*User) to avoid overwrite immutable fields
	UpdateUser(ctx context.Context, user *User) error
}
