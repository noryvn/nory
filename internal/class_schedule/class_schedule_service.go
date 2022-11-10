package classschedule

import (
	"context"

	"nory/common/response"
	"nory/common/validator"
	"nory/domain"
)

type ClassScheduleService struct {
	ClassScheduleRepository domain.ClassScheduleRepository
}

func (css *ClassScheduleService) CreateSchedule(ctx context.Context, schedule *domain.ClassSchedule) (*response.Response[any], error) {
	if err := validator.ValidateStruct(schedule); err != nil {
		return nil, err
	}
	if err := css.ClassScheduleRepository.CreateSchedule(ctx, schedule); err != nil {
		return nil, err
	}
	return response.New[any](204, nil), nil
}
