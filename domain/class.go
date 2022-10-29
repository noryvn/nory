package domain

import (
	"context"
	"errors"
	"time"
)

var (
	// class not found
	ErrClassNotExists = errors.New("class does not exists")
	// conflicting data, unique constraint and etc.
	ErrClassAlreadyExists = errors.New("class already exists")
)

type Class struct {
	ClassId   string    `json:"classId"`                 // immutable, unique
	OwnerId   string    `json:"ownerId" validate:"uuid"` // immutable
	CreatedAt time.Time `json:"createdAt"`               // immutable

	Name        string `json:"name" validate:"required,max=20"` // mutable
	Description string `json:"description" validate:"max=255"`  // mutable
}

func (c *Class) Update(cc *Class) {
	if cc.Name != "" {
		c.Name = cc.Name
	}
	if cc.Description != "" {
		c.Description = cc.Description
	}
}

type ClassRepository interface {
	GetClass(ctx context.Context, classId string) (*Class, error)
	GetClassesByOwnerId(ctx context.Context, ownerId string) ([]*Class, error)
	CreateClass(ctx context.Context, class *Class) error
	DeleteClass(ctx context.Context, classId string) error
	UpdateClass(ctx context.Context, class *Class) error
}
