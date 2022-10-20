package domain

import "context"

type ClassSchedule struct {
	ScheduleId string `json:"scheduleId"`
	ClassId    string `json:"-"`

	Name     string `json:"name"`
	StartAt  string `json:"startAt"`
	Duration int16  `json:"duration"`
	Day      int8   `json:"day"`
}

type ClassScheduleRepository interface {
	CreateSchedule(ctx context.Context, schedule *ClassSchedule) error
	GetSchedules(ctx context.Context, classId string) ([]ClassSchedule, error)
	DeleteSchedule(ctx context.Context, scheduleId string) error
}
