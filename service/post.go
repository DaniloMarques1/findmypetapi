package service

import (
	"github.com/danilomarques1/findmypetapi/dto"
	"github.com/danilomarques1/findmypetapi/model"
	"github.com/google/uuid"
)

type PostService struct {
	postRepository model.PostRepository
}

func NewPostService(postRepository model.PostRepository) *PostService {
	return &PostService{
		postRepository: postRepository,
	}
}

func (ps *PostService) CreatePost(postDto dto.CreatePostRequestDto,
	userId string) (*dto.CreatePostResponseDto, error) {
	// TODO produce a rabbit mq message

	postId := uuid.NewString()
	post := model.Post{
		Id:          postId,
		AuthorId:    userId,
		Title:       postDto.Title,
		Description: postDto.Description,
		ImageUrl:    "/path/to/image",
	}
	err := ps.postRepository.Save(&post)
	if err != nil {
		return nil, err
	}

	response := dto.CreatePostResponseDto{
		Post: post,
	}

	return &response, nil
}

func (ps *PostService) GetAll() (*dto.GetPostResponseDto, error) {
	posts, err := ps.postRepository.FindAll()
	if err != nil {
		return nil, err
	}

	response := dto.GetPostResponseDto{Posts: posts}
	return &response, nil
}
