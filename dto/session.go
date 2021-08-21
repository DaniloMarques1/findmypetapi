package dto

import "github.com/danilomarques1/findmypetapi/model"

type SessionRequestDto struct {
	Email    string `json:"email" validate:"required,max=60,email"`
	Password string `json:"password" validate:"required",max=20"`
}

type SessionResponseDto struct {
	Token        string     `json:"token"`
	RefreshToken string     `json:"refresh_token"`
	User         model.User `json:"user"`
}
