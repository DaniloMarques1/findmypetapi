package service

import (
	"log"
	"net/http"

	"github.com/danilomarques1/findmypetapi/dto"
	"github.com/danilomarques1/findmypetapi/model"
	"github.com/danilomarques1/findmypetapi/util"
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
		return nil, util.NewApiError("Password and confirm password does not match", http.StatusBadRequest)
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	id := uuid.NewString()
	user := model.User{
		Id:           id,
		Name:         userDto.Name,
		Email:        userDto.Email,
		PasswordHash: string(passwordHash),
	}

	userAr, err := us.userRepo.FindByEmail(userDto.Email);
	if  err == nil {
		if userAr != nil && userAr.Email == userDto.Email {
			log.Printf("Email already registered %v\n", err)
			return nil, util.NewApiError("Email already taken", http.StatusBadRequest)
		}
	}


	if err := us.userRepo.Save(&user); err != nil {
		return nil, err
	}

	return &dto.CreateUserResponseDto{
		User: user,
	}, nil
}