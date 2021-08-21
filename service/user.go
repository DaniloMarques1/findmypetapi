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

	userAr, err := us.userRepo.FindByEmail(userDto.Email)
	if err == nil {
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

func (us *UserService) CreateSession(sessionRequest dto.SessionRequestDto) (
	*dto.SessionResponseDto, error) {
	user, err := us.userRepo.FindByEmail(sessionRequest.Email)
	if err != nil {
		return nil, util.NewApiError("Invalid email", http.StatusBadRequest)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(sessionRequest.Password)); err != nil {
		return nil, util.NewApiError("Invalid password", http.StatusBadRequest)
	}

	token, refreshToken, err := util.NewToken(user.Id)
	if err != nil {
		return nil, err
	}
	response := dto.SessionResponseDto{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
	}

	return &response, nil
}

func (us *UserService) RefreshSession(refreshToken string) (*dto.SessionResponseDto, error) {
	userClaims, err := util.VerifyToken(refreshToken)
	if err != nil {
		return nil, util.NewApiError("Invalid token", http.StatusUnauthorized)
	}

	if userClaims.TokenRole != util.REFRESH_TOKEN {
		return nil, util.NewApiError("The given token is not a refresh token",
			http.StatusUnauthorized)
	}

	user, err := us.userRepo.FindById(userClaims.UserId)
	if err != nil {
		return nil, util.NewApiError("User was not found", http.StatusNotFound)
	}

	token, refreshToken, err := util.NewToken(user.Id)
	if err != nil {
		return nil, err
	}
	response := dto.SessionResponseDto{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
	}

	return &response, nil
}
