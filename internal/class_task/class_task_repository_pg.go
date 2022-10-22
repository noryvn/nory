package classtask

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"

	"nory/domain"
)

type ClassTaskRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewClassTaskRepositoryPostgres(pool *pgxpool.Pool) *ClassTaskRepositoryPostgres {
	return &ClassTaskRepositoryPostgres{pool}
}

func (ctrp *ClassTaskRepositoryPostgres) CreateTask(ctx context.Context, task *domain.ClassTask) error {
	task.TaskId = xid.New().String()
	_, err := ctrp.pool.Exec(
		ctx,
		"INSERT INTO class_task(task_id, class_id, name, description, due_date) VALUES($1, $2, $3, $4, $5);",
		task.TaskId,
		task.ClassId,
		task.Name,
		task.Description,
		task.DueDate,
	)
	return err
}

func (ctrp *ClassTaskRepositoryPostgres) GetTask(ctx context.Context, taskId string) (*domain.ClassTask, error) {
	ct := &domain.ClassTask{
		TaskId: taskId,
	}
	row := ctrp.pool.QueryRow(
		ctx,
		"SELECT class_id, name, description, due_date FROM class_task WHERE task_id = $1",
		taskId,
	)
	err := row.Scan(
		&ct.ClassId,
		&ct.Name,
		&ct.Description,
		&ct.DueDate,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		err = domain.ErrClassTaskNotExists
	}
	return ct, err
}

func (ctrp *ClassTaskRepositoryPostgres) GetTasks(ctx context.Context, classId string) ([]*domain.ClassTask, error) {
	return ctrp.GetTasksWithRange(ctx, classId, time.Unix(1, 0), time.Date(2030, time.August, 11, 0, 0, 0, 0, time.UTC))
}

func (ctrp *ClassTaskRepositoryPostgres) GetTasksWithRange(ctx context.Context, classId string, from, to time.Time) ([]*domain.ClassTask, error) {
	var tasks []*domain.ClassTask
	rows, err := ctrp.pool.Query(
		ctx,
		"SELECT task_id, name, description, due_date FROM class_task WHERE class_id = $1 AND due_date >= $2 AND due_date < $3",
		classId,
		from,
		to,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		ct := &domain.ClassTask{
			ClassId: classId,
		}
		err := rows.Scan(
			&ct.TaskId,
			&ct.Name,
			&ct.Description,
			&ct.DueDate,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, ct)
	}
	return tasks, nil
}

func (ctrp *ClassTaskRepositoryPostgres) UpdateTask(ctx context.Context, task *domain.ClassTask) error {
	return nil
}

func (ctrp *ClassTaskRepositoryPostgres) DeleteTask(ctx context.Context, taskId string) error {
	_, err := ctrp.pool.Exec(
		ctx,
		"DELETE FROM class_task WHERE task_id = $1",
		taskId,
	)
	return err
}
