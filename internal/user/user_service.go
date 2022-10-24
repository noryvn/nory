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
