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

type PostHandler struct {
	postService *service.PostService
	validator   *validator.Validate
}

func NewPostHandler(postService *service.PostService, validator *validator.Validate) *PostHandler {
	return &PostHandler{
		postService: postService,
		validator:   validator,
	}
}

func (ph *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var postDto dto.CreatePostRequestDto
	if err := json.NewDecoder(r.Body).Decode(&postDto); err != nil {
		log.Printf("Invalid body %v\n", err)
		util.RespondJson(w, http.StatusBadRequest, dto.ErrorDto{Message: "Invalid body"})
		return
	}
	if err := ph.validator.Struct(postDto); err != nil {
		log.Printf("Invalid body %v\n", err)
		util.RespondJson(w, http.StatusBadRequest, dto.ErrorDto{Message: "Invalid body"})
		return
	}
	userId := r.Header.Get("user_id")

	response, err := ph.postService.CreatePost(postDto, userId)
	if err != nil {
		log.Printf("Error creating post %v\n", err)
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusCreated, response)
}
