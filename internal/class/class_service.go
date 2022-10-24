package class

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nory/common/response"
	"nory/common/validator"
	"nory/domain"
)

type ClassService struct {
	ClassRepository     domain.ClassRepository
	ClassTaskRepository domain.ClassTaskRepository
}

func (cs *ClassService) GetClassInfo(ctx context.Context, classId string) (*response.Response[*domain.Class], error) {
	class, err := cs.ClassRepository.GetClass(ctx, classId)
	if errors.Is(err, domain.ErrClassNotExists) {
		msg := fmt.Sprintf("can not find class with id %q", classId)
		return nil, response.NewNotFound(msg)
	}
	if err != nil {
		return nil, err
	}
	return response.New(200, class), err
}

func (cs *ClassService) GetClassTasks(ctx context.Context, classId string, from, to time.Time) (*response.Response[[]*domain.ClassTask], error) {
	if from.IsZero() {
		from = time.Now()
	}
	if to.IsZero() {
		to = from.Add(24 * time.Hour)
	}
	tasks, err := cs.ClassTaskRepository.GetTasksWithRange(ctx, classId, from, to)
	return response.New(200, tasks), err
}

func (cs *ClassService) CreateClass(ctx context.Context, class *domain.Class) (*response.Response[*domain.Class], error) {
	if err := validator.ValidateStruct(class); err != nil {
		return nil, err
	}
	err := cs.ClassRepository.CreateClass(ctx, class)
	return response.New(200, class), err
}

func (cs *ClassService) CreateClassTask(ctx context.Context, task *domain.ClassTask) (*response.Response[*domain.ClassTask], error) {
	if err := validator.ValidateStruct(task); err != nil {
		return nil, err
	}
	err := cs.ClassTaskRepository.CreateTask(ctx, task)
	return response.New(200, task), err
}

func (cs *ClassService) DeleteClass(ctx context.Context, classId, userId string) (*response.Response[any], error) {
	class, err := cs.ClassRepository.GetClass(ctx, classId)
	if err != nil {
		return nil, err
	}

	if class.OwnerId != userId {
		msg := fmt.Sprintf("user with id %q does not has access to class with id %q", userId, classId)
		return nil, response.NewForbidden(msg)
	}

	if err := cs.ClassRepository.DeleteClass(ctx, classId); err != nil {
		return nil, err
	}

	return response.New[any](204, nil), nil
}
