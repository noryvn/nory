package user

import (
	"context"
	"errors"
	"fmt"

	"nory/common/response"
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

func (us UserService) GetUserProfile(ctx context.Context, user *domain.User) (*response.Response[*UserProfile], error) {
	user, err := us.UserRepository.GetUser(ctx, user.UserId)
	if err != nil {
		return nil, err
	}

	classes, err := us.ClassRepository.GetClassesByOwnerId(ctx, user.UserId)
	if err != nil {
		return nil, err
	}

	members, err := us.ClassMemberRepository.ListJoined(ctx, user.UserId)
	if err != nil {
		return nil, err
	}

	up := &UserProfile{
		User:        user,
		OwnedClass:  len(classes),
		JoinedClass: len(members),
	}
	return response.New(200, up), nil
}

func (us UserService) GetUserProfileById(ctx context.Context, userId string) (*response.Response[*UserProfile], error) {
	user, err := us.UserRepository.GetUser(ctx, userId)
	if errors.Is(err, domain.ErrUserNotExists) {
		msg := fmt.Sprintf("can not find user with id %q", userId)
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
