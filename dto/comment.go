package dto

import "github.com/danilomarques1/findmypetapi/model"

type CreateCommentRequestDto struct {
	CommentText string `json:"comment_text" validate:"required,max=400"`
}

type CreateCommentResponseDto struct {
	Comment model.Comment `json:"comment"`
}
