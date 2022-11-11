package class

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nory/common/response"
	"nory/common/validator"
	"nory/domain"
)

type ClassService struct {
	UserRepository          domain.UserRepository
	ClassRepository         domain.ClassRepository
	ClassTaskRepository     domain.ClassTaskRepository
	ClassMemberRepository   domain.ClassMemberRepository
	ClassScheduleRepository domain.ClassScheduleRepository
}

func (cs *ClassService) GetClassInfo(ctx context.Context, classId string) (*response.Response[*domain.Class], error) {
	class, err := cs.ClassRepository.GetClass(ctx, classId)
	if errors.Is(err, domain.ErrClassNotExists) {
		msg := fmt.Sprintf("can not find class with id %q", classId)
		return nil, response.NewNotFound(msg)
	}
	if err != nil {
		return nil, err
	}
	return response.New(200, class), nil
}

func (cs *ClassService) GetClassInfoByName(ctx context.Context, ownerId, name string)  (*response.Response[*domain.Class], error) {
	class, err := cs.ClassRepository.GetClassByName(ctx, ownerId, name)
	if errors.Is(err, domain.ErrClassNotExists) {
		msg := fmt.Sprintf("can not find class with name %q that owned by %q", name, ownerId)
		return nil, response.NewNotFound(msg)
	}
	if err != nil {
		return nil, err
	}
	return response.New(200, class), nil
}

func (cs *ClassService) GetClassTasks(ctx context.Context, classId string, from, to time.Time) (*response.Response[[]*domain.ClassTask], error) {
	if from.IsZero() {
		from = time.Now()
	}
	if to.IsZero() {
		to = from.Add(7 * 24 * time.Hour)
	}
	tasks, err := cs.ClassTaskRepository.GetTasksWithRange(ctx, classId, from, to)
	return response.New(200, tasks), err
}

func (cs *ClassService) CreateClass(ctx context.Context, class *domain.Class) (*response.Response[*domain.Class], error) {
	if err := validator.ValidateStruct(class); err != nil {
		return nil, err
	}
	if err := cs.ClassRepository.CreateClass(ctx, class); err != nil {
		return nil, err
	}
	if err := cs.ClassMemberRepository.CreateMember(ctx, &domain.ClassMember{
		UserId:  class.OwnerId,
		ClassId: class.ClassId,
		Level:   "owner",
	}); err != nil {
		return nil, err
	}
	return response.New(200, class), nil
}

func (cs *ClassService) CreateClassTask(ctx context.Context, userId string, task *domain.ClassTask) (*response.Response[*domain.ClassTask], error) {
	if err := validator.ValidateStruct(task); err != nil {
		return nil, err
	}
	if err := cs.AccessClass(ctx, userId, task.ClassId, "member"); err != nil {
		return nil, err
	}
	if err := cs.ClassTaskRepository.CreateTask(ctx, task); err != nil {
		return nil, err
	}
	return response.New(200, task), nil
}

func (cs *ClassService) DeleteClassTask(ctx context.Context, userId, taskId string) (*response.Response[any], error) {
	task, err := cs.ClassTaskRepository.GetTask(ctx, taskId)
	if errors.Is(err, domain.ErrClassTaskNotExists) {
		msg := fmt.Sprintf("can not find task with id %q", taskId)
		return nil, response.NewUnprocessableEntity(msg)
	}
	if err != nil {
		return nil, err
	}

	if err := cs.AccessClass(ctx, userId, task.ClassId, "admin"); err != nil {
		return nil, err
	}

	if err := cs.ClassTaskRepository.DeleteTask(ctx, taskId); err != nil {
		return nil, err
	}

	return response.New[any](204, nil), nil
}

func (cs *ClassService) AddMember(ctx context.Context, userId string, member *domain.ClassMember) (*response.Response[any], error) {
	if err := validator.ValidateStruct(member); err != nil {
		return nil, err
	}
	if err := cs.AccessClass(ctx, userId, member.ClassId, "admin"); err != nil {
		return nil, err
	}
	if err := cs.ClassMemberRepository.CreateMember(ctx, member); err != nil {
		return nil, err
	}
	return response.New[any](204, nil), nil
}

func (cs *ClassService) AddMemberByUsername(ctx context.Context, userId, username string, member *domain.ClassMember) (*response.Response[any], error) {
	user, err := cs.UserRepository.GetUserByUsername(ctx, username)
	if errors.Is(err, domain.ErrUserNotExists) {
		msg := fmt.Sprintf("can not find user with username %q", username)
		return nil, response.NewUnprocessableEntity(msg)
	}
	if err != nil {
		return nil, err
	}
	member.UserId = user.UserId
	return cs.AddMember(ctx, userId, member)
}

func (cs *ClassService) DeleteMember(ctx context.Context, userId, classId, memberId string) (*response.Response[any], error) {
	if err := cs.AccessClass(ctx, userId, classId, "admin"); err != nil {
		return nil, err
	}
	if err := cs.ClassMemberRepository.DeleteMember(ctx, &domain.ClassMember{ClassId: classId, UserId: memberId}); err != nil {
		return nil, err
	}
	return response.New[any](204, nil), nil
}

func (cs *ClassService) ListMember(ctx context.Context, classId string) (*response.Response[[]*domain.ClassMember], error) {
	members, err := cs.ClassMemberRepository.ListMembers(ctx, classId)
	if err != nil {
		return nil, err
	}

	return response.New(200, members), nil
}

func (cs *ClassService) UpdateMember(ctx context.Context, userId string, member *domain.ClassMember) (*response.Response[any], error) {
	if err := cs.AccessClass(ctx, userId, member.ClassId, "admin"); err != nil {
		return nil, err
	}
	if err := validator.ValidateStruct(member); err != nil {
		return nil, err
	}
	if err := cs.ClassMemberRepository.UpdateMember(ctx, member); err != nil {
		return nil, err
	}
	return response.New[any](204, nil), nil
}

func (cs *ClassService) DeleteClass(ctx context.Context, userId, classId string) (*response.Response[any], error) {
	if err := cs.AccessClass(ctx, userId, classId, "admin"); err != nil {
		return nil, err
	}

	if err := cs.ClassRepository.DeleteClass(ctx, classId); err != nil {
		return nil, err
	}

	return response.New[any](204, nil), nil
}

func (cs *ClassService) CreateSchedule(ctx context.Context, schedule *domain.ClassSchedule) (*response.Response[any], error) {
	if err := cs.AccessClass(ctx, schedule.AuthorId, schedule.ClassId, "admin"); err != nil {
		return nil, err
	}
	if err := cs.ClassScheduleRepository.CreateSchedule(ctx, schedule); err != nil {
		return nil, err
	}
	return response.New[any](204, nil), nil
}

func (cs *ClassService) DeleteSchedule(ctx context.Context, userId, scheduleId string) (*response.Response[any], error) {
	schedule, err := cs.ClassScheduleRepository.GetSchedule(ctx, scheduleId)
	if errors.Is(err, domain.ErrClassScheduleNotExists) {
		msg := fmt.Sprintf("can not find class schedule with id %q", scheduleId)
		return nil, response.NewUnprocessableEntity(msg)
	}
	if err != nil {
		return nil, err
	}

	if err := cs.AccessClass(ctx, userId, schedule.ClassId, "admin"); err != nil {
		return nil, err
	}
	if err := cs.ClassScheduleRepository.DeleteSchedule(ctx, scheduleId); err != nil {
		return nil, err
	}
	return response.New[any](204, nil), nil
}

func (cs *ClassService) ClearSchedules(ctx context.Context, userId, classId string, day int8) (*response.Response[any], error) {
	if err := cs.AccessClass(ctx, userId, classId, "admin"); err != nil {
		return nil, err
	}
	if err := cs.ClassScheduleRepository.ClearSchedules(ctx, classId, day); err != nil {
		return nil, err
	}
	return response.New[any](204, nil), nil
}

func (cs *ClassService) GetClassSchedules(ctx context.Context, classId string) (*response.Response[[]*domain.ClassSchedule], error) {
	schedules, err := cs.ClassScheduleRepository.GetSchedules(ctx, classId)
	if err != nil {
		return nil, err
	}
	return response.New(200, schedules), nil
}

func (cs *ClassService) GetSchedule(ctx context.Context, scheduleId string) (*response.Response[*domain.ClassSchedule], error) {
	schedules, err := cs.ClassScheduleRepository.GetSchedule(ctx, scheduleId)
	if errors.Is(err, domain.ErrClassScheduleNotExists) {
		msg := fmt.Sprintf("can not find class schedule with id %q", scheduleId)
		return nil, response.NewNotFound(msg)
	}
	if err != nil {
		return nil, err
	}
	return response.New(200, schedules), nil
}

type permissionLevel uint8

const (
	permissionOwner permissionLevel = 255 - iota
	permissionAdmin
	permissionUser
)

func permissionLevelFromString(s string) permissionLevel {
	switch s {
	case "owner": return permissionOwner
	case "admin": return permissionAdmin
	default: return permissionUser
	}
}

func (cs *ClassService) AccessClass(ctx context.Context, userId, classId, minimum string) error {
	msg := fmt.Sprintf("user with id %q does not has %q access to class with id %q", userId, minimum, classId)
	resErr := response.NewForbidden(msg)

	member, err := cs.ClassMemberRepository.GetMember(ctx, &domain.ClassMember{
		ClassId: classId,
		UserId:  userId,
	})
	if errors.Is(err, domain.ErrClassMemberNotExists) {
		return resErr
	}
	if err != nil {
		return err
	}

	required := permissionLevelFromString(minimum)
	memberLevel := permissionLevelFromString(member.Level)

	if memberLevel >= required {
		return nil
	}

	return resErr
}
