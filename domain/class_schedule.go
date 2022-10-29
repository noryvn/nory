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
	ScheduleId string    `json:"scheduleId"` // unique
	ClassId    string    `json:"classId"`    //
	AuthorId   string    `json:"authorId"`
	CreatedAt  time.Time `json:"createdAt"`

	Name     string        `json:"name"`     //
	StartAt  time.Duration `json:"startAt"`  //
	Duration int16         `json:"duration"` //
	Day      int8          `json:"day"`      //
}

type ClassScheduleRepository interface {
	CreateSchedule(ctx context.Context, schedule *ClassSchedule) error
	GetSchedule(ctx context.Context, scheduleId string) (*ClassSchedule, error)
	GetSchedules(ctx context.Context, classId string) ([]*ClassSchedule, error)
	DeleteSchedule(ctx context.Context, scheduleId string) error
	ClearSchedules(ctx context.Context, classId string, day int8) error
}
