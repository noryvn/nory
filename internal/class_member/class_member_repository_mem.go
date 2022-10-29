package classmember

import (
	"context"
	"sync"

	"nory/domain"
)

type ClassMemberRepositoryMem struct {
	mx sync.Mutex
	members []*domain.ClassMember
}

func NewClassMemberRepositoryMem() *ClassMemberRepositoryMem {
	return &ClassMemberRepositoryMem{}
}

func (repo *ClassMemberRepositoryMem) ListMembers(ctx context.Context, classId string) ([]*domain.ClassMember, error) {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	var result []*domain.ClassMember
	for _, m := range repo.members {
		m := m
		if m.ClassId == classId {
			result = append(result, m)
		}
	}
	return result, nil
}

func (repo *ClassMemberRepositoryMem) ListJoined(ctx context.Context, userId string) ([]*domain.ClassMember, error) {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	var result []*domain.ClassMember
	for _, m := range repo.members {
		m := m
		if m.UserId == userId {
			result = append(result, m)
		}
	}
	return result, nil
}

func (repo *ClassMemberRepositoryMem) GetMember(ctx context.Context, member *domain.ClassMember) (*domain.ClassMember, error) {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	for _, m := range repo.members {
		m := m
		if m.ClassId == member.ClassId && m.UserId == member.UserId {
			return m, nil
		}
	}
	return nil, domain.ErrClassMemberNotExists
}

func (repo *ClassMemberRepositoryMem) CreateMember(ctx context.Context, member *domain.ClassMember) error {
	_, err := repo.GetMember(ctx, member)
	if err == nil {
		return domain.ErrClassMemberAlreadyExists
	}

	repo.mx.Lock()
	defer repo.mx.Unlock()

	repo.members = append(repo.members, member)
	return nil
}

func (repo *ClassMemberRepositoryMem) UpdateMember(ctx context.Context, member *domain.ClassMember) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	for _, m := range repo.members {
		m := m
		if m.ClassId == member.ClassId && m.UserId == member.UserId {
			m.Update(member)
		}
	}
	return nil
}

func (repo *ClassMemberRepositoryMem) DeleteMember(ctx context.Context, member *domain.ClassMember) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	for i, m := range repo.members {
		if m.ClassId == member.ClassId && m.UserId == member.UserId {
			repo.members = append(repo.members[:i], repo.members[i+1:]...)
			return nil
		}
	}
	return nil
}
