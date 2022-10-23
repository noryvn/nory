package user

import (
	"context"

	"nory/common/response"
	"nory/domain"
)

type UserProfile struct {
	domain.User

	OwnedClass int `json:"ownedClass"`
}

func (up *UserProfile) GetUser() *domain.User {
	u := &domain.User{}
	u.UserId = up.UserId
	u.CreatedAt = up.CreatedAt
	u.Username = up.Username
	u.Name = up.Name
	u.Email = up.Email
	return u
}

func (up *UserProfile) SetUser(u *domain.User) {
	up.UserId = u.UserId
	up.CreatedAt = u.CreatedAt
	up.Username = u.Username
	up.Name = u.Name
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
	up.SetUser(user)
	up.OwnedClass = len(classes)
	return res, nil
}
