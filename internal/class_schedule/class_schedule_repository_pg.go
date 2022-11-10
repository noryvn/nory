package classschedule

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"nory/domain"
)

type ClassScheduleRepositoryPg struct {
	pool *pgxpool.Pool
}

func NewClassScheduleRepositoryPg(pool *pgxpool.Pool) *ClassScheduleRepositoryPg {
	return &ClassScheduleRepositoryPg{pool}
}

func (csrp *ClassScheduleRepositoryPg) CreateSchedule(ctx context.Context, schedule *domain.ClassSchedule) error {
	return nil
}

func (csrp *ClassScheduleRepositoryPg) GetSchedule(ctx context.Context, scheduleId string) (*domain.ClassSchedule, error) {
	return nil, nil
}

func (csrp *ClassScheduleRepositoryPg) GetSchedules(ctx context.Context, classId string) ([]*domain.ClassSchedule, error) {
	return nil, nil
}

func (csrp *ClassScheduleRepositoryPg) DeleteSchedule(ctx context.Context, scheduleId string) error {
	return nil
}

func (csrp *ClassScheduleRepositoryPg) ClearSchedules(ctx context.Context, classId string, day int8) error {
	return nil
}

