package user

import (
	"context"

	"nory/common/response"
	"nory/domain"
)

type UserProfile struct {
	User *domain.User `json:"user"`

	OwnedClass int `json:"ownedClass"`
}

type UserService struct {
	UserRepository  domain.UserRepository
	ClassRepository domain.ClassRepository
}

func (us UserService) GetUserProfile(ctx context.Context, user *domain.User) (*response.Response[*UserProfile], error) {
	up := &UserProfile{}
	res := response.New(200, up)
	classes, err := us.ClassRepository.GetClassesByOwnerId(ctx, user.UserId)
	if err != nil {
		return res, err
	}
	up.User = user
	up.OwnedClass = len(classes)
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
