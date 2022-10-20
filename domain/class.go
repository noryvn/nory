package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrClassNotFound = errors.New("class not found")
	ErrClassExists   = errors.New("class already exists")
)

type Class struct {
	ClassId   string    `json:"classId"`
	OwnerId   string    `json:"ownerId"`
	CreatedAt time.Time `json:"createdAt"`

	Name        string `json:"name"`
	Description string `json:"description"`
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
	GetByOwnerId(ctx context.Context, ownerId string) ([]*Class, error)
	CreateClass(ctx context.Context, class *Class) error
	DeleteClass(ctx context.Context, classId string) error
	UpdateClass(ctx context.Context, class *Class) error
}
