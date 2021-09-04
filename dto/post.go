package dto

import "github.com/danilomarques1/findmypetapi/model"

type CreatePostRequestDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreatePostResponseDto struct {
	Post model.Post `json:"post"`
}

type GetPostResponseDto struct {
	Posts []model.Post `json:"posts"`
}
