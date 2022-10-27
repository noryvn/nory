package classtask

import (
	"context"
	"fmt"

	"nory/common/response"
	"nory/common/validator"
	"nory/domain"
)

type ClassTaskService struct {
	ClassRepository     domain.ClassRepository
	ClassTaskRepository domain.ClassTaskRepository
}

func (cts *ClassTaskService) CreateTask(ctx context.Context, task *domain.ClassTask) (*response.Response[*domain.ClassTask], error) {
	if err := validator.ValidateStruct(task); err != nil {
		return nil, err
	}
	if err := cts.ClassTaskRepository.CreateTask(ctx, task); err != nil {
		return nil, err
	}
	return response.New(200, task), nil
}

func (cts *ClassTaskService) GetTask(ctx context.Context, taskId string) (*response.Response[*domain.ClassTask], error) {
	task, err := cts.ClassTaskRepository.GetTask(ctx, taskId)
	if err != nil {
		return nil, err
	}

	return response.New(200, task), nil
}

func (cts *ClassTaskService) DeleteTask(ctx context.Context, taskId, userId string) (*response.Response[*struct{}], error) {
	task, err := cts.ClassTaskRepository.GetTask(ctx, taskId)
	if err != nil {
		return nil, err
	}

	class, err := cts.ClassRepository.GetClass(ctx, task.ClassId)
	if err != nil {
		return nil, err
	}

	if class.OwnerId != userId {
		msg := fmt.Sprintf("user with id %q does not has access to class with id %q", userId, task.ClassId)
		return nil, response.NewForbidden(msg)
	}

	if err := cts.ClassTaskRepository.DeleteTask(ctx, task.TaskId); err != nil {
		return nil, err
	}

	return response.New[*struct{}](204, nil), nil
}
