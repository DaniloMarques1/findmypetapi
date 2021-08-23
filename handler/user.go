package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/danilomarques1/findmypetapi/dto"
	"github.com/danilomarques1/findmypetapi/service"
	"github.com/danilomarques1/findmypetapi/util"
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
		log.Printf("Error decoding body %v\n", err)
		util.RespondJson(w, http.StatusBadRequest, dto.ErrorDto{Message: "Invalid body"})
		return
	}

	if err := uh.validator.Struct(userDto); err != nil {
		log.Printf("Error validating struct %v\n", err)
		util.RespondJson(w, http.StatusBadRequest, dto.ErrorDto{Message: "Invalid body"})
		return
	}

	response, err := uh.userService.Save(userDto)
	if err != nil {
		log.Printf("Error saving user %v\n", err)
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusCreated, response)
}

func (uh *UserHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var sessionRequest dto.SessionRequestDto
	if err := json.NewDecoder(r.Body).Decode(&sessionRequest); err != nil {
		util.RespondJson(w, http.StatusBadRequest, dto.ErrorDto{"Invalid body"})
		return
	}

	if err := uh.validator.Struct(sessionRequest); err != nil {
		util.RespondJson(w, http.StatusBadRequest, dto.ErrorDto{"Invalid body"})
		return
	}

	response, err := uh.userService.CreateSession(sessionRequest)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusOK, response)
}

func (uh *UserHandler) RefreshSession(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("refresh_token")
	response, err := uh.userService.RefreshSession(refreshToken)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusOK, response)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("user_id")
	log.Printf("%v\n", userId)
	var updateDto dto.UpdateUserRequestDto
	if err := json.NewDecoder(r.Body).Decode(&updateDto); err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorDto{Message: "Invalid body"})
		return
	}
	if err := uh.validator.Struct(updateDto); err != nil {
		util.RespondJson(w, http.StatusBadRequest, &dto.ErrorDto{Message: "Invalid body"})
		return
	}

	err := uh.userService.UpdateUser(userId, updateDto)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusNoContent, nil)
}
