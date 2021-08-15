package model

import (
	"github.com/google/uuid"
)

type User struct {
	id 	     uuid.UUID
	Name         string
	Email        string
	PasswordHash string
}

type UserRepository interface {
	Save(*User) error
	FindById(id uuid.UUID) (*User, error)
	Update(*User) error
	FindByEmail(email string) (*User, error)
}
