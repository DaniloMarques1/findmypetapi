package handler

import (
	"encoding/json"
	"net/http"

	"github.com/danilomarques1/findmypetapi/dto"
	"github.com/danilomarques1/findmypetapi/service"
	"github.com/danilomarques1/findmypetapi/util"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type CommentHandler struct {
	commentService *service.CommentService
	validator      *validator.Validate
}

func NewCommentHandler(commentService *service.CommentService,
	validator *validator.Validate) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		validator:      validator,
	}
}

func (ch *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var createComment dto.CreateCommentRequestDto
	if err := json.NewDecoder(r.Body).Decode(&createComment); err != nil {
		util.RespondJson(w, http.StatusBadRequest,
			dto.ErrorDto{Message: "Invalid body"})
		return
	}
	if err := ch.validator.Struct(createComment); err != nil {
		util.RespondJson(w, http.StatusBadRequest,
			dto.ErrorDto{Message: "Invalid body"})
		return
	}
	vars := mux.Vars(r)
	postId := vars["post_id"]
	userId := r.Header.Get("user_id")
	response, err := ch.commentService.Save(userId, postId, createComment)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.RespondJson(w, http.StatusCreated, response)
}
