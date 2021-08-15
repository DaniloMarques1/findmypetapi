package model

type User struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type UserRepository interface {
	Save(*User) error
	FindById(id string) (*User, error)
	Update(*User) error
	FindByEmail(email string) (*User, error)
}
