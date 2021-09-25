package dto

import (
	"time"

	"github.com/danilomarques1/findmypetapi/model"
)

type CreateCommentRequestDto struct {
	CommentText string `json:"comment_text" validate:"required,max=400"`
}

type CreateCommentResponseDto struct {
	Comment CommentDto `json:"comment"`
}

type CommentDto struct {
	Id          string    `json:"id"`
	AuthorId    string    `json:"author_id"`
	PostId      string    `json:"post_id"`
	CommentText string    `json:"comment_text"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetCommentsResponseDto struct {
	Comments []model.GetCommentDto `json:"comments"`
}
