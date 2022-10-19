package user

import (
	"context"
	"sync"

	"nory/domain"
)

type UserRepositoryMem struct {
	mu sync.Mutex
	m  map[string]*domain.User
}

func NewUserRepositoryMem() *UserRepositoryMem {
	return &UserRepositoryMem{
		m: make(map[string]*domain.User),
	}
}

func (urm *UserRepositoryMem) GetUser(ctx context.Context, id string) (*domain.User, error) {
	urm.mu.Lock()
	defer urm.mu.Unlock()
	u, ok := urm.m[id]
	if !ok {
		return u, domain.ErrUserNotFound
	}
	return u, nil
}

func (urm *UserRepositoryMem) CreateUser(ctx context.Context, u *domain.User) error {
	urm.mu.Lock()
	defer urm.mu.Unlock()
	if _, ok := urm.m[u.UserId]; ok {
		return domain.ErrUserExists
	}
	for _, user := range urm.m {
		if user.Username == u.Username {
			return domain.ErrUserExists
		}
	}
	urm.m[u.UserId] = u
	return nil
}

func (urm *UserRepositoryMem) DeleteUser(ctx context.Context, id string) error {
	urm.mu.Lock()
	defer urm.mu.Unlock()
	delete(urm.m, id)
	return nil
}

func (urm *UserRepositoryMem) UpdateUser(ctx context.Context, u *domain.User) error {
	for _, uu := range urm.m {
		if uu.UserId == u.UserId {
			continue
		}
		if uu.Username == u.Username {
			return domain.ErrUserExists
		}
	}
	uu, err := urm.GetUser(ctx, u.UserId)
	if err != nil {
		return err
	}
	if u.Name != "" {
		uu.Name = u.Name
	}
	if u.Username != "" {
		uu.Username = u.Username
	}
	return nil
}
