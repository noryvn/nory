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

func (urm *UserRepositoryMem) GetUserByUserId(ctx context.Context, id string) (*domain.User, error) {
	urm.mu.Lock()
	defer urm.mu.Unlock()
	u, ok := urm.m[id]
	if !ok {
		return u, domain.ErrUserNotExists
	}
	return u, nil
}

func (urm *UserRepositoryMem) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	urm.mu.Lock()
	defer urm.mu.Unlock()
	for _, user := range urm.m {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotExists
}

func (urm *UserRepositoryMem) CreateUser(ctx context.Context, u *domain.User) error {
	urm.mu.Lock()
	defer urm.mu.Unlock()
	if _, ok := urm.m[u.UserId]; ok {
		return domain.ErrUserAlreadyExists
	}
	for _, user := range urm.m {
		if user.Username == u.Username || user.Email == u.Email {
			return domain.ErrUserAlreadyExists
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
		if uu.Username == u.Username || uu.Email == u.Email {
			return domain.ErrUserAlreadyExists
		}
	}
	uu, err := urm.GetUserByUserId(ctx, u.UserId)
	if err != nil {
		return err
	}
	uu.Update(u)
	return nil
}
