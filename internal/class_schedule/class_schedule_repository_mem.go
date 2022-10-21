package classschedule

import (
	"context"
	"nory/domain"
	"sync"

	"github.com/rs/xid"
)

type ClassScheduleRepositoryMem struct {
	mx sync.Mutex
	m  map[string]*domain.ClassSchedule
}

func NewClassScheduleRepositoryMem() *ClassScheduleRepositoryMem {
	return &ClassScheduleRepositoryMem{
		m: make(map[string]*domain.ClassSchedule),
	}
}

func (csrm *ClassScheduleRepositoryMem) CreateSchedule(ctx context.Context, schedule *domain.ClassSchedule) error {
	csrm.mx.Lock()
	defer csrm.mx.Unlock()
	schedule.ScheduleId = xid.New().String()
	csrm.m[schedule.ScheduleId] = schedule
	return nil
}

func (csrm *ClassScheduleRepositoryMem) GetSchedule(ctx context.Context, scheduleId string) (*domain.ClassSchedule, error) {
	csrm.mx.Lock()
	defer csrm.mx.Unlock()
	schedule, ok := csrm.m[scheduleId]
	if !ok {
		return nil, domain.ErrClassScheduleNotFound
	}
	return schedule, nil
}

func (csrm *ClassScheduleRepositoryMem) GetSchedules(ctx context.Context, classId string) ([]*domain.ClassSchedule, error) {
	csrm.mx.Lock()
	defer csrm.mx.Unlock()
	schedules := make([]*domain.ClassSchedule, 0)
	for _, sch := range csrm.m {
		if sch.ClassId == classId {
			schedules = append(schedules, sch)
		}
	}
	return schedules, nil
}

func (csrm *ClassScheduleRepositoryMem) DeleteSchedule(ctx context.Context, scheduleId string) error {
	csrm.mx.Lock()
	defer csrm.mx.Unlock()
	delete(csrm.m, scheduleId)
	return nil
}

func (csrm *ClassScheduleRepositoryMem) ClearSchedules(ctx context.Context, classId string, day int8) error {
	csrm.mx.Lock()
	defer csrm.mx.Unlock()
	for key, schedule := range csrm.m {
		if schedule.ClassId == classId && schedule.Day == day {
			delete(csrm.m, key)
		}
	}
	return nil
}
