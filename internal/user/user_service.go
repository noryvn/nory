package user

import (
	"context"
	"errors"
	"fmt"

	"nory/common/response"
	"nory/common/validator"
	"nory/domain"
)

type UserProfile struct {
	User *domain.User `json:"user"`

	OwnedClass  int `json:"ownedClass"`
	JoinedClass int `json:"joinedClass"`
}

type UserService struct {
	UserRepository        domain.UserRepository
	ClassRepository       domain.ClassRepository
	ClassMemberRepository domain.ClassMemberRepository
}

func (us UserService) GetUserProfile(ctx context.Context, user *domain.User) (*response.Response[*domain.User], error) {
	classes, err := us.ClassRepository.GetClassesByOwnerId(ctx, user.UserId)
	if err != nil {
		return nil, err
	}

	members, err := us.ClassMemberRepository.ListJoined(ctx, user.UserId)
	if err != nil {
		return nil, err
	}

	user.UserStatistics = &domain.UserStatistics{
		OwnedClass:  len(classes),
		JoinedClass: len(members),
	}
	user.OwnedClass = classes

	return response.New(200, user), nil
}

func (us UserService) GetUserProfileById(ctx context.Context, userId string) (*response.Response[*domain.User], error) {
	user, err := us.UserRepository.GetUserByUserId(ctx, userId)
	if errors.Is(err, domain.ErrUserNotExists) {
		msg := fmt.Sprintf("can not find user with id %q", userId)
		return nil, response.NewNotFound(msg)
	}
	if err != nil {
		return nil, err
	}
	return us.GetUserProfile(ctx, user)
}

func (us UserService) GetUserProfileByUsername(ctx context.Context, username string) (*response.Response[*domain.User], error) {
	user, err := us.UserRepository.GetUserByUsername(ctx, username)
	if errors.Is(err, domain.ErrUserNotExists) {
		msg := fmt.Sprintf("can not find user with id %q", username)
		return nil, response.NewNotFound(msg)
	}
	if err != nil {
		return nil, err
	}
	return us.GetUserProfile(ctx, user)
}

func (us UserService) GetUserClasses(ctx context.Context, user *domain.User) (*response.Response[[]*domain.Class], error) {
	classes, err := us.ClassRepository.GetClassesByOwnerId(ctx, user.UserId)
	return response.New(200, classes), err
}

func (us UserService) GetUserJoinedClasses(ctx context.Context, user *domain.User) (*response.Response[[]*domain.ClassMember], error) {
	classes, err := us.ClassMemberRepository.ListJoined(ctx, user.UserId)
	if err != nil {
		return nil, err
	}
	return response.New(200, classes), nil
}

func (us UserService) UpdateUser(ctx context.Context, user *domain.User) (*response.Response[any], error) {
	if err := validator.ValidateStruct(user); err != nil {
		return nil, err
	}
	if err := us.UserRepository.UpdateUser(ctx, user); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return nil, response.NewConflict("user already exists")
		}
		return nil, err
	}
	return response.New[any](204, nil), nil
}
