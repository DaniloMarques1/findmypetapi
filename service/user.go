package service

import (
	"github.com/danilomarques1/findmypetapi/dto"
	"github.com/danilomarques1/findmypetapi/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo model.UserRepository
}

func NewUserService(userRepo model.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (us *UserService) Save(userDto dto.CreateUserRequestDto) (*dto.CreateUserResponseDto, error) {
	if userDto.Password != userDto.ConfirmPassword {
		// TODO error handling
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.MinCost)
	if err != nil {
		// TODO
	}
	id := uuid.NewString()
	user := model.User{
		Id:           id,
		Name:         userDto.Name,
		Email:        userDto.Email,
		PasswordHash: string(passwordHash),
	}

	if err := us.userRepo.Save(&user); err != nil {
		// TODO
	}

	return &dto.CreateUserResponseDto{
		User: user,
	}, nil
}
