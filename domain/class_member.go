package domain

import (
	"context"
	"errors"
)

var (
	ErrClassMemberNotExists = errors.New("ClassMember does not exists")
)

type ClassMember struct {
	ClassId string
	UserId string

	Level string
}

func (member *ClassMember) Update(m *ClassMember) {
	if m.Level != "" {
		member.Level = m.Level
	}
}

type ClassMemberRepository interface {
	ListMembers(ctx context.Context, classId string) ([]*ClassMember, error)
	ListJoined(ctx context.Context, userId string) ([]*ClassMember, error)
	GetMember(ctx context.Context, member *ClassMember) (*ClassMember, error)
	CreateMember(ctx context.Context, member *ClassMember) error
	UpdateMember(ctx context.Context, member *ClassMember) error
	DeleteMember(ctx context.Context, member *ClassMember) error
}
