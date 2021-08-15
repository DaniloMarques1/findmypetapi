package service

import (
	"github.com/danilomarques1/findmypetapi/model"
)

type UserService struct {
	userRepo model.UserRepository
}

func NewUserService(userRepo model.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}
