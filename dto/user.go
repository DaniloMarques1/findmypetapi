package dto

import (
	"github.com/danilomarques1/findmypetapi/model"
)

type CreateUserRequestDto struct {
	Name            string `json:"name" validate:"required,max=100"`
	Email           string `json:"email" validate:"required,max=60,email"`
	Password        string `json:"password" validate:"required,max=20"`
	ConfirmPassword string `json:"confirm_password" validate:"required,max=20"`
}

type CreateUserResponseDto struct {
	User model.User `json:"user"`
}

type UpdateUserDto struct {
	Name            string `json:"name" validate:"required,max=100"`
	OldPassword     string `json:"old_password" validate:"required,max=20"`
	NewPassword     string `json:"new_password" validate:"required,max=20"`
	ConfirmPassword string `json:"confirm_password" validate:"required,max=20"`
}
