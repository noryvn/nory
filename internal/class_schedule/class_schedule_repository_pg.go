package classschedule

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"

	"nory/domain"
)

type ClassScheduleRepositoryPg struct {
	pool *pgxpool.Pool
}

func NewClassScheduleRepositoryPg(pool *pgxpool.Pool) *ClassScheduleRepositoryPg {
	return &ClassScheduleRepositoryPg{pool}
}

func (csrp *ClassScheduleRepositoryPg) CreateSchedule(ctx context.Context, schedule *domain.ClassSchedule) error {
	schedule.ScheduleId = xid.New().String()

	_, err := csrp.pool.Exec(
		ctx,
		`INSERT INTO class_schedule(schedule_id, class_id, author_id, name, start_at, duration, day) VALUES($1, $2, $3, $4, $5, $6)`,
		schedule.ScheduleId,
		schedule.ClassId,
		schedule.AuthorId,
		schedule.Name,
		schedule.Duration,
		schedule.Day,
	)

	return err
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
	_, err := csrp.pool.Exec(
		ctx,
		"DELETE FROM class_schedule WHERE day = $1",
		day,
	)
	return err
}

