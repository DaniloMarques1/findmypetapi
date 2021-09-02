package dto

import "github.com/danilomarques1/findmypetapi/model"

type CreatePostRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreatePostResponse struct {
	Post model.Post `json:"post"`
}

type GetPostResponse struct {
	Post []model.Post `json:"posts"`
}
