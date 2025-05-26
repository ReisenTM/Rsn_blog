package user_service

import "blogX_server/model"

type UserService struct {
	userModel model.UserModel
}

func NewUserService(user model.UserModel) *UserService {
	return &UserService{
		userModel: user,
	}
}
