package handler

import (
	"encoding/json"
	"net/http"

	"github.com/danilomarques1/findmypetapi/dto"
	"github.com/danilomarques1/findmypetapi/util"
	"github.com/danilomarques1/findmypetapi/service"
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

	response, err := uh.userService.Save(userDto)
	if err != nil {
		// TODO error handling
	}

	util.RespondJson(w, http.StatusCreated, response)
}
