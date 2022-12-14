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
		"INSERT INTO class_task(task_id, class_id, author_id, author_display_name, name, description, due_date) VALUES($1, $2, $3, $4, $5, $6, $7);",
		task.TaskId,
		task.ClassId,
		task.AuthorId,
		task.AuthorDisplayName,
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
		"SELECT class_id, author_id, created_at, author_display_name, name, description, due_date FROM class_task WHERE task_id = $1 ORDER BY due_date",
		taskId,
	)
	err := row.Scan(
		&ct.ClassId,
		&ct.AuthorId,
		&ct.CreatedAt,
		&ct.AuthorDisplayName,
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
	tasks := make([]*domain.ClassTask, 0)
	rows, err := ctrp.pool.Query(
		ctx,
		"SELECT task_id, author_id, created_at, author_display_name, name, description, due_date FROM class_task WHERE class_id = $1 AND due_date >= $2 AND due_date < $3 ORDER BY task_id",
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
			&ct.AuthorId,
			&ct.CreatedAt,
			&ct.AuthorDisplayName,
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
	ct, err := ctrp.GetTask(ctx, task.TaskId)
	if err != nil {
		return err
	}
	ct.Update(task)
	_, err = ctrp.pool.Exec(
		ctx,
		"UPDATE class_task SET name = $1, description = $2, due_date = $3 WHERE task_id = $4",
		ct.Name,
		ct.Description,
		ct.DueDate,
		ct.TaskId,
	)
	return err
}

func (ctrp *ClassTaskRepositoryPostgres) DeleteTask(ctx context.Context, taskId string) error {
	_, err := ctrp.pool.Exec(
		ctx,
		"DELETE FROM class_task WHERE task_id = $1",
		taskId,
	)
	return err
}
