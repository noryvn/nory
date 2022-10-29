package class

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"

	"nory/domain"
)

type ClassRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewClassRepositoryPostgres(pool *pgxpool.Pool) *ClassRepositoryPostgres {
	return &ClassRepositoryPostgres{pool}
}

func (crp *ClassRepositoryPostgres) GetClass(ctx context.Context, classId string) (*domain.Class, error) {
	class := &domain.Class{
		ClassId: classId,
	}
	row := crp.pool.QueryRow(ctx, "SELECT owner_id, created_at, name, description FROM class WHERE class_id = $1", classId)
	err := row.Scan(
		&class.OwnerId,
		&class.CreatedAt,
		&class.Name,
		&class.Description,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		err = domain.ErrClassNotExists
	}
	return class, err
}

func (crp *ClassRepositoryPostgres) GetClassesByOwnerId(ctx context.Context, ownerId string) ([]*domain.Class, error) {
	classes := make([]*domain.Class, 0)
	rows, err := crp.pool.Query(ctx, "SELECT class_id, created_at, name, description FROM class WHERE owner_id = $1 ORDER BY class_id", ownerId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		class := &domain.Class{OwnerId: ownerId}
		if err := rows.Scan(
			&class.ClassId,
			&class.CreatedAt,
			&class.Name,
			&class.Description,
		); err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}
	return classes, nil
}

func (crp *ClassRepositoryPostgres) CreateClass(ctx context.Context, class *domain.Class) error {
	class.ClassId = xid.New().String()
	_, err := crp.pool.Exec(
		ctx,
		"INSERT INTO class(class_id, owner_id, name, description) VALUES($1, $2, $3, $4)",
		class.ClassId,
		class.OwnerId,
		class.Name,
		class.Description,
	)
	return err
}

func (crp *ClassRepositoryPostgres) DeleteClass(ctx context.Context, classId string) error {
	_, err := crp.pool.Exec(
		ctx,
		"DELETE FROM class WHERE class_id = $1",
		classId,
	)
	return err
}

func (crp *ClassRepositoryPostgres) UpdateClass(ctx context.Context, class *domain.Class) error {
	c, err := crp.GetClass(ctx, class.ClassId)
	if err != nil {
		return err
	}
	c.Update(class)
	_, err = crp.pool.Exec(
		ctx,
		"UPDATE class SET name = $1, description = $2 WHERE class_id = $3",
		c.Name,
		c.Description,
		c.ClassId,
	)
	return err
}
