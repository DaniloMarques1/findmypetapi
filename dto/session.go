package dto

import "github.com/danilomarques1/findmypetapi/model"

type SessionRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SessionResponseDto struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}
