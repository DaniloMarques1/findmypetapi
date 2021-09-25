package model

import (
	"time"
)

type Comment struct {
	Id          string    `json:"id"`
	AuthorId    string    `json:"author_id"`
	PostId      string    `json:"post_id"`
	CommentText string    `json:"comment_text"`
	CreatedAt   time.Time `json:"created_at"`
}

// TODO should be on dto
type GetCommentDto struct {
	Id          string    `json:"id"`
	CommentText string    `json:"comment_text"`
	CreatedAt   time.Time `json:"created_at"`
	Author      AuthorDto    `json:"author"`
	PostId      string    `json:"post_id"`
}

// TODO should be on dto
type AuthorDto struct {
	AuthorName  string `json:"author_name"`
	AuthorEmail string `json:"author_email"`
}

type CommentRepository interface {
	Save(*Comment) error
	FindAll(postId string) ([]GetCommentDto, error)
	GetCommentNotificationMessage(string, string) ([]byte, error)
}
