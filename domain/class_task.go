package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrClassTaskNotExists = errors.New("class task does not exists")
)

type ClassTask struct {
	TaskId    string    `json:"taskId"`    // immutable, unique
	ClassId   string    `json:"classId"`   // immutable
	AuthorId  string    `json:"authorId"`  // immutable
	CreatedAt time.Time `json:"createdAt"` // immutable

	AuthorDisplayName string    `json:"authorDisplayName" validate:"max=20"` // mutable
	Name              string    `json:"name" validate:"max=20"`              // mutable
	Description       string    `json:"description" validate:"max=1024"`     // mutable
	DueDate           time.Time `json:"dueDate"`                             // mutable
}

func (ct *ClassTask) Update(task *ClassTask) {
	if task.Name != "" {
		ct.Name = task.Name
	}
	if task.Description != "" {
		ct.Description = task.Description
	}
	if !task.DueDate.IsZero() {
		ct.DueDate = task.DueDate
	}
}

type ClassTaskRepository interface {
	// CreateTask should update (*ClassTask).TaskId to generated id from database or etc.
	CreateTask(ctx context.Context, task *ClassTask) error
	GetTask(ctx context.Context, taskId string) (*ClassTask, error)
	GetTasks(ctx context.Context, classId string) ([]*ClassTask, error)
	GetTasksWithRange(ctx context.Context, classId string, from, to time.Time) ([]*ClassTask, error)
	UpdateTask(ctx context.Context, task *ClassTask) error
	DeleteTask(ctx context.Context, taskId string) error
}
