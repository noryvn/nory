package class

import (
	"context"
	"sync"

	"nory/domain"

	"github.com/rs/xid"
)

type ClassRepositoryMem struct {
	mx sync.Mutex
	m  map[string]*domain.Class
}

func NewClassRepositoryMem() *ClassRepositoryMem {
	return &ClassRepositoryMem{
		m: make(map[string]*domain.Class),
	}
}

func (crm *ClassRepositoryMem) GetClass(ctx context.Context, classId string) (*domain.Class, error) {
	crm.mx.Lock()
	defer crm.mx.Unlock()
	c, ok := crm.m[classId]
	if !ok {
		return nil, domain.ErrClassNotFound
	}
	return c, nil
}

func (crm *ClassRepositoryMem) GetClassesByOwnerId(ctx context.Context, ownerId string) ([]*domain.Class, error) {
	crm.mx.Lock()
	defer crm.mx.Unlock()
	var classes []*domain.Class
	for _, c := range crm.m {
		if c.OwnerId != ownerId {
			continue
		}
		classes = append(classes, c)
	}
	return classes, nil
}

func (crm *ClassRepositoryMem) CreateClass(ctx context.Context, class *domain.Class) error {
	crm.mx.Lock()
	defer crm.mx.Unlock()
	class.ClassId = xid.New().String()
	crm.m[class.ClassId] = class
	return nil
}

func (crm *ClassRepositoryMem) DeleteClass(ctx context.Context, classId string) error {
	crm.mx.Lock()
	defer crm.mx.Unlock()
	delete(crm.m, classId)
	return nil
}

func (crm *ClassRepositoryMem) UpdateClass(ctx context.Context, class *domain.Class) error {
	c, err := crm.GetClass(ctx, class.ClassId)
	if err != nil {
		return err
	}
	if class.Description != "" {
		c.Description = class.Description
	}
	if class.Name != "" {
		c.Name = class.Name
	}
	return nil
}
