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
	ClassId   string
	OwnerId   string
	CreatedAt time.Time

	Name        string
	Description string
}

type ClassRepository interface {
	GetClass(ctx context.Context, classId string) (*Class, error)
	GetByOwnedBy(ctx context.Context, ownerId string) ([]*Class, error)
	CreateClass(ctx context.Context, class *Class) error
	DeleteClass(ctx context.Context, classId string) error
	UpdateClass(ctx context.Context, class *Class) error
}
