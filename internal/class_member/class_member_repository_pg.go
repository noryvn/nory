package classmember

import (
	"context"
	"errors"
	"nory/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClassMemberRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewClassMemberRepositoryPostgres(pool *pgxpool.Pool) *ClassMemberRepositoryPostgres {
	return &ClassMemberRepositoryPostgres{pool}
}

func (repo *ClassMemberRepositoryPostgres) ListMembers(ctx context.Context, classId string) ([]*domain.ClassMember, error) {
	members := make([]*domain.ClassMember, 0)
	rows, err := repo.pool.Query(
		ctx,
		"SELECT user_id, created_at, level FROM class_member WHERE class_id = $1 ORDER BY created_at",
		classId,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		member := &domain.ClassMember{
			ClassId: classId,
		}

		err := rows.Scan(
			&member.UserId,
			&member.CreatedAt,
			&member.Level,
		)

		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}
	return members, nil
}

func (repo *ClassMemberRepositoryPostgres) ListJoined(ctx context.Context, userId string) ([]*domain.ClassMember, error) {
	members := make([]*domain.ClassMember, 0)
	rows, err := repo.pool.Query(
		ctx,
		"SELECT class_id, created_at, level FROM class_member WHERE user_id = $1 ORDER BY created_at",
		userId,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		member := &domain.ClassMember{
			UserId: userId,
		}

		err := rows.Scan(
			&member.ClassId,
			&member.CreatedAt,
			&member.Level,
		)

		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}
	return members, nil
}

func (repo *ClassMemberRepositoryPostgres) GetMember(ctx context.Context, member *domain.ClassMember) (*domain.ClassMember, error) {
	m := &domain.ClassMember{
		UserId:  member.UserId,
		ClassId: member.ClassId,
	}
	row := repo.pool.QueryRow(
		ctx,
		"SELECT level, created_at FROM class_member WHERE user_id = $1 AND class_id = $2",
		member.UserId,
		member.ClassId,
	)
	err := row.Scan(
		&m.Level,
		&m.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrClassMemberNotExists
	}
	return m, nil
}

func (repo *ClassMemberRepositoryPostgres) CreateMember(ctx context.Context, member *domain.ClassMember) error {
	_, err := repo.pool.Exec(
		ctx,
		"INSERT INTO class_member(class_id, user_id, level) VALUES($1, $2, $3)",
		member.ClassId,
		member.UserId,
		member.Level,
	)
	if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
		return domain.ErrClassMemberAlreadyExists
	}
	return err
}

func (repo *ClassMemberRepositoryPostgres) UpdateMember(ctx context.Context, member *domain.ClassMember) error {
	_, err := repo.pool.Exec(
		ctx,
		"UPDATE class_member SET level = $1 WHERE user_id = $2 AND class_id = $3",
		member.Level,
		member.UserId,
		member.ClassId,
	)
	return err
}

func (repo *ClassMemberRepositoryPostgres) DeleteMember(ctx context.Context, member *domain.ClassMember) error {
	_, err := repo.pool.Exec(
		ctx,
		"DELETE FROM class_member WHERE user_id = $1 AND class_id = $2",
		member.UserId,
		member.ClassId,
	)
	return err
}
