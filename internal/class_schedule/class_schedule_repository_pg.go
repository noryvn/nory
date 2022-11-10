package classschedule

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
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
		`INSERT INTO class_schedule(schedule_id, class_id, author_id, name, start_at, duration, day) VALUES($1, $2, $3, $4, $5, $6, $7)`,
		schedule.ScheduleId,
		schedule.ClassId,
		schedule.AuthorId,
		schedule.Name,
		schedule.StartAt,
		schedule.Duration,
		schedule.Day,
	)

	return err
}

func (csrp *ClassScheduleRepositoryPg) GetSchedule(ctx context.Context, scheduleId string) (*domain.ClassSchedule, error) {
	schedule := &domain.ClassSchedule{
		ScheduleId: scheduleId,
	}
	row := csrp.pool.QueryRow(
		ctx,
		"SELECT class_id, author_id, created_at, name, start_at, duration, day FROM class_schedule WHERE schedule_id = $1",
		scheduleId,
	)
	err := row.Scan(
		&schedule.ClassId,
		&schedule.AuthorId,
		&schedule.CreatedAt,
		&schedule.Name,
		&schedule.StartAt,
		&schedule.Duration,
		&schedule.Day,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		err = domain.ErrClassScheduleNotExists
	}
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

func (csrp *ClassScheduleRepositoryPg) GetSchedules(ctx context.Context, classId string) ([]*domain.ClassSchedule, error) {
	var schedules []*domain.ClassSchedule

	rows, err := csrp.pool.Query(
		ctx,
		"SELECT schedule_id, author_id, created_at, name, start_at, duration, day FROM class_schedule WHERE class_id = $1 ORDER BY schedule_id",
		classId,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		schedule := &domain.ClassSchedule{
			ClassId: classId,
		}

		err := rows.Scan(
			&schedule.ScheduleId,
			&schedule.AuthorId,
			&schedule.CreatedAt,
			&schedule.Name,
			&schedule.StartAt,
			&schedule.Duration,
			&schedule.Day,
		)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (csrp *ClassScheduleRepositoryPg) DeleteSchedule(ctx context.Context, scheduleId string) error {
	_, err := csrp.pool.Exec(
		ctx,
		"DELETE FROM class_schedule WHERE schedule_id = $1",
		scheduleId,
	)
	return err
}

func (csrp *ClassScheduleRepositoryPg) ClearSchedules(ctx context.Context, classId string, day int8) error {
	_, err := csrp.pool.Exec(
		ctx,
		"DELETE FROM class_schedule WHERE day = $1",
		day,
	)
	return err
}
