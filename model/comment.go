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

type CommentRepository interface {
	Save(*Comment) error
	FindAll(postId string) ([]Comment, error)
}
