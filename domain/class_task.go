package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrClassTaskNotExists = errors.New("class not found")
)

type ClassTask struct {
	TaskId  string `json:"taskId"`
	ClassId string `json:"-"`

	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
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
	// CreateTask should update (ClassTask).TaskId to generated id from database or etc.
	CreateTask(ctx context.Context, task *ClassTask) error
	GetTask(ctx context.Context, taskId string) (*ClassTask, error)
	GetTasks(ctx context.Context, classId string) ([]*ClassTask, error)
	GetTasksWithRange(ctx context.Context, classId string, from, to time.Time) ([]*ClassTask, error)
	UpdateTask(ctx context.Context, task *ClassTask) error
	DeleteTask(ctx context.Context, taskId string) error
}
