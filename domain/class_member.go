package domain

import "context"

type ClassMember struct {
	ClassId string
	UserId string

	Level string
}

type ClassMemberRepository interface {
	ListMembers(ctx context.Context, classId string) ([]*ClassMember, error)
	IsMember(ctx context.Context, classId, memberId string) (bool, error)
	CreateMember(ctx context.Context, member *ClassMember) error
	DeleteMember(ctx context.Context, classId, memberId string) error
}
