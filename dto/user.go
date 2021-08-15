package dto

import (
	"github.com/danilomarques1/findmypetapi/model"
)

type CreateUserRequestDto struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type CreateUserResponseDto struct {
	User model.User `json:"user"`
}
