package class

import (
	"context"

	"nory/common/response"
	"nory/domain"
)

type ClassService struct {
	ClassRepository     domain.ClassRepository
	ClassTaskRepository domain.ClassTaskRepository
}

func (cs *ClassService) GetClassInfo(ctx context.Context, classId string) (*response.Response[*domain.Class], error) {
	class, err := cs.ClassRepository.GetClass(ctx, classId)
	return response.New(200, class), err
}

func (cs *ClassService) GetClassTasks(ctx context.Context, classId string) (*response.Response[[]*domain.ClassTask], error) {
	return nil, nil
}
