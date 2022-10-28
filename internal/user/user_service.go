package user

import (
	"context"

	"nory/common/response"
	"nory/domain"
)

type UserProfile struct {
	User *domain.User `json:"user"`

	OwnedClass int `json:"ownedClass"`
	JoinedClass int `json:"joinedClass"`
}

type UserService struct {
	UserRepository  domain.UserRepository
	ClassRepository domain.ClassRepository
	ClassMemberRepository domain.ClassMemberRepository
}

func (us UserService) GetUserProfile(ctx context.Context, user *domain.User) (*response.Response[*UserProfile], error) {
	up := &UserProfile{}
	res := response.New(200, up)
	classes, err := us.ClassRepository.GetClassesByOwnerId(ctx, user.UserId)
	if err != nil {
		return res, err
	}

	members, err := us.ClassMemberRepository.ListJoined(ctx, user.UserId)
	up.User = user
	up.OwnedClass = len(classes)
	up.JoinedClass = len(members)
	return res, nil
}

func (us UserService) GetUserProfileById(ctx context.Context, userId string) (*response.Response[*UserProfile], error) {
	user, err := us.UserRepository.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return us.GetUserProfile(ctx, user)
}

func (us UserService) GetUserClasses(ctx context.Context, user *domain.User) (*response.Response[[]*domain.Class], error) {
	classes, err := us.ClassRepository.GetClassesByOwnerId(ctx, user.UserId)
	return response.New(200, classes), err
}
