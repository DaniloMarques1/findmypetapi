package dto

import "github.com/danilomarques1/findmypetapi/model"

type CreatePostRequestDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreatePostResponseDto struct {
	Post model.Post `json:"post"`
}

type GetPostsResponseDto struct {
	Posts []model.Post `json:"posts"`
}

type GetPostResponseDto struct {
	Post model.Post `json:"post"`
}
