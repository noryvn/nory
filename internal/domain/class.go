package domain

import (
	"context"
	"time"
)

type Class struct {
	ClassId   string
	OwnerId   string
	CreatedAt time.Time

	Name        string
	Description string
}

type ClassRepository interface {
	Create(ctx *context.Context, class Class) error
	GetById(ctx *context.Context, classId string) (Class, error)
	GetOwnedBy(ctx *context.Context, ownerId string) ([]Class, error)
	Delete(ctx *context.Context, classId string) error
}
