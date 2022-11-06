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
	UserRepository        domain.UserRepository
	ClassRepository       domain.ClassRepository
	ClassTaskRepository   domain.ClassTaskRepository
	ClassMemberRepository domain.ClassMemberRepository
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
	if err := cs.AccessClass(ctx, userId, task.ClassId); err != nil {
		return nil, err
	}
	if err := cs.ClassTaskRepository.CreateTask(ctx, task); err != nil {
		return nil, err
	}
	return response.New(200, task), nil
}

func (cs *ClassService) AddMember(ctx context.Context, userId string, member *domain.ClassMember) (*response.Response[any], error) {
	if err := validator.ValidateStruct(member); err != nil {
		return nil, err
	}
	if err := cs.AccessClass(ctx, userId, member.ClassId); err != nil {
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
	if err := cs.AccessClass(ctx, userId, classId); err != nil {
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
	if err := cs.AccessClass(ctx, userId, member.ClassId); err != nil {
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
	if err := cs.AccessClass(ctx, userId, classId); err != nil {
		return nil, err
	}

	if err := cs.ClassRepository.DeleteClass(ctx, classId); err != nil {
		return nil, err
	}

	return response.New[any](204, nil), nil
}

func (cs *ClassService) AccessClass(ctx context.Context, userId, classId string) error {
	class, err := cs.ClassRepository.GetClass(ctx, classId)
	if errors.Is(err, domain.ErrClassNotExists) {
		msg := fmt.Sprintf("can not find class with id %q", classId)
		return response.NewNotFound(msg)
	}
	if err != nil {
		return err
	}

	if class.OwnerId == userId {
		return nil
	}

	member, err := cs.ClassMemberRepository.GetMember(ctx, &domain.ClassMember{
		ClassId: classId,
		UserId:  userId,
	})

	if err == nil && member.Level == "admin" {
		return nil
	}

	msg := fmt.Sprintf("user with id %q does not has modify access to class with id %q", userId, classId)
	return response.NewForbidden(msg)
}
