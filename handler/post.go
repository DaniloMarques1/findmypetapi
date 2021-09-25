package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/danilomarques1/findmypetapi/dto"
	"github.com/danilomarques1/findmypetapi/service"
	"github.com/danilomarques1/findmypetapi/util"
	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
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

const MAX_BYTES_SIZE_POST_CREATION = 10000

func (ph *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var postDto dto.CreatePostRequestDto
	r.ParseMultipartForm(MAX_BYTES_SIZE_POST_CREATION)

	postDto.Title = r.FormValue("title")
	postDto.Description = r.FormValue("description")

	file, fileHandler, err := r.FormFile("file")
	if err != nil {
		util.RespondJson(w, http.StatusBadRequest, dto.ErrorDto{Message: "The given file is invalid"})
		return
	}
	defer file.Close()

	postDto.File = file
	postDto.Filename = fmt.Sprintf("%v-%v", time.Now().Unix(), fileHandler.Filename)
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

func (ph *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	response, err := ph.postService.GetAll()
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusOK, response)
}

func (ph *PostHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId := vars["post_id"]
	response, err := ph.postService.FindById(postId)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusOK, response)
}

func (ph *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	var updateDto dto.UpdatePostRequestDto
	if err := json.NewDecoder(r.Body).Decode(&updateDto); err != nil {
		log.Printf("Error parsing body %v\n", err)
		util.RespondJson(w, http.StatusBadRequest,
			dto.ErrorDto{Message: "Invalid body"})
		return
	}
	if err := ph.validator.Struct(updateDto); err != nil {
		log.Printf("Error validating body %v\n", err)
		util.RespondJson(w, http.StatusBadRequest,
			dto.ErrorDto{Message: "Invalid body"})
		return
	}

	vars := mux.Vars(r)
	postId := vars["post_id"]
	userId := r.Header.Get("user_id")
	err := ph.postService.Update(updateDto, userId, postId)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusNoContent, nil)
}

func (ph *PostHandler) FindPostsByAuthor(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("user_id")
	response, err := ph.postService.FindPostsByAuthor(userId)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusOK, response)
}
