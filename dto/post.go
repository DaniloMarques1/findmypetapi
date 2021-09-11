package dto

import "github.com/danilomarques1/findmypetapi/model"

type CreatePostRequestDto struct {
	Title       string `json:"title" validate:"max=120"`
	Description string `json:"description" validate:"max=800"`
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

type UpdatePostRequestDto struct {
	Title       string `json:"title" validate:"required,max=120"`
	Description string `json:"description" validate:"required,max=800"`
	Status      string `json:"status" validate:"required"`
}
