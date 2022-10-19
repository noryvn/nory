package class

import "nory/internal/domain"

type ClassService struct {
	classRepository domain.ClassRepository
}

func NewClassService(
	cr domain.ClassRepository,
) *ClassService {
	return &ClassService{
		cr,
	}
}
