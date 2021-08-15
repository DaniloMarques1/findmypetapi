package handler

import (
	"encoding/json"
	"net/http"

	"github.com/danilomarques1/findmypetapi/service"
	"github.com/danilomarques1/findmypetapi/dto"
	validator "github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService *service.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService *service.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

func (uh *UserHandler) Save(w http.ResponseWriter, r *http.Request) {
	var userDto dto.CreateUserRequestDto
	if err := json.NewDecoder(r.Body).Decode(&userDto); err != nil {
		// TODO error handling
	}

	if err := uh.validator.Struct(userDto); err != nil {
		// TODO error handling
	}

	_, err := uh.userService.Save(userDto) // TODO get and return response
	if err != nil {
		// TODO error handling
	}

	// TODO return 201 and response

}
