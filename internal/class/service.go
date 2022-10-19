package class

import "nory/domain"

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
