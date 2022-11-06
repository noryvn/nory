package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrClassScheduleNotExists = errors.New("class schedule does not exists")
)

type ClassSchedule struct {
	ScheduleId string    `json:"scheduleId"`  // immutable, unique
	ClassId    string    `json:"classId"`     // immutable
	AuthorId   string    `json:"authorId"` // immutable
	CreatedAt  time.Time `json:"createdAt"` // immutable

	Name     string        `json:"name"`     // immutable
	StartAt  time.Duration `json:"startAt"`  // immutable
	Duration int16         `json:"duration"` // immutable
	Day      int8          `json:"day"`      // immutable
}

type ClassScheduleRepository interface {
	CreateSchedule(ctx context.Context, schedule *ClassSchedule) error
	GetSchedule(ctx context.Context, scheduleId string) (*ClassSchedule, error)
	GetSchedules(ctx context.Context, classId string) ([]*ClassSchedule, error)
	DeleteSchedule(ctx context.Context, scheduleId string) error
	ClearSchedules(ctx context.Context, classId string, day int8) error
}
