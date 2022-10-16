package user

import (
	"context"
	"sync"

	"nory/internal/domain"
)

type RepositoryMem struct {
	mu sync.Mutex
	m  map[string]*domain.User
}

func NewRepositoryMem() *RepositoryMem {
	return &RepositoryMem{
		m: make(map[string]*domain.User),
	}
}

func (rm *RepositoryMem) GetUser(ctx context.Context, id string) (*domain.User, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	u, ok := rm.m[id]
	if !ok {
		return u, domain.ErrUserNotFound
	}
	return u, nil
}

func (rm *RepositoryMem) CreateUser(ctx context.Context, u *domain.User) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if _, ok := rm.m[u.UserId]; ok {
		return domain.ErrDuplicateUser
	}
	for _, user := range rm.m {
		if user.Username == u.Username {
			return domain.ErrDuplicateUser
		}
	}
	rm.m[u.UserId] = u
	return nil
}

func (rm *RepositoryMem) DeleteUser(ctx context.Context, id string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.m, id)
	return nil
}

func (rm *RepositoryMem) UpdateUser(ctx context.Context, u *domain.User) error {
	for _, uu := range rm.m {
		if uu.UserId == u.UserId {
			continue
		}
		if uu.Username == u.Username {
			return domain.ErrDuplicateUser
		}
	}
	uu, err := rm.GetUser(ctx, u.UserId)
	if err != nil {
		return err
	}
	uu.Name = u.Name
	uu.Username = u.Username
	return nil
}
