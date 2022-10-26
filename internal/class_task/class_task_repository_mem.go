package classtask

import (
	"context"
	"errors"
	"sync"
	"time"

	"nory/domain"

	"github.com/rs/xid"
)

type ClassTaskRepositoryMem struct {
	mx sync.Mutex
	m  map[string]*domain.ClassTask
}

func NewClassTaskRepositoryMem() *ClassTaskRepositoryMem {
	return &ClassTaskRepositoryMem{
		m: make(map[string]*domain.ClassTask),
	}
}

func (ctrm *ClassTaskRepositoryMem) CreateTask(ctx context.Context, task *domain.ClassTask) error {
	if task.ClassId == "" || task.AuthorId == "" {
		return errors.New("empty data")
	}
	ctrm.mx.Lock()
	defer ctrm.mx.Unlock()
	task.TaskId = xid.New().String()
	ctrm.m[task.TaskId] = task
	return nil
}

func (ctrm *ClassTaskRepositoryMem) GetTask(ctx context.Context, taskId string) (*domain.ClassTask, error) {
	ctrm.mx.Lock()
	defer ctrm.mx.Unlock()
	task, ok := ctrm.m[taskId]
	if !ok {
		return nil, domain.ErrClassTaskNotExists
	}
	return task, nil
}

func (ctrm *ClassTaskRepositoryMem) GetTasks(ctx context.Context, classId string) ([]*domain.ClassTask, error) {
	return ctrm.GetTasksWithRange(ctx, classId, time.Unix(1, 0), time.Date(2030, time.August, 11, 0, 0, 0, 0, time.UTC))
}

func (ctrm *ClassTaskRepositoryMem) GetTasksWithRange(ctx context.Context, classId string, from, to time.Time) ([]*domain.ClassTask, error) {
	ctrm.mx.Lock()
	defer ctrm.mx.Unlock()
	tasks := make([]*domain.ClassTask, 0)
	for _, task := range ctrm.m {
		if task.ClassId == classId &&
			!task.DueDate.Before(from) &&
			task.DueDate.Before(to) {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (ctrm *ClassTaskRepositoryMem) GetTasksWithDate(ctx context.Context, classId string, dueDate time.Time) ([]*domain.ClassTask, error) {
	ctrm.mx.Lock()
	defer ctrm.mx.Unlock()
	tasks := make([]*domain.ClassTask, 0)
	for _, task := range ctrm.m {
		if task.ClassId == classId && dueDate.Equal(task.DueDate) {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (ctrm *ClassTaskRepositoryMem) UpdateTask(ctx context.Context, task *domain.ClassTask) error {
	ctrm.mx.Lock()
	defer ctrm.mx.Unlock()
	t, ok := ctrm.m[task.TaskId]
	if !ok {
		return domain.ErrClassTaskNotExists
	}
	t.Update(task)
	return nil
}

func (ctrm *ClassTaskRepositoryMem) DeleteTask(ctx context.Context, taskId string) error {
	ctrm.mx.Lock()
	defer ctrm.mx.Unlock()
	delete(ctrm.m, taskId)
	return nil
}
