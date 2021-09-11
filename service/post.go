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

func (ps *PostService) GetAll() (*dto.GetPostsResponseDto, error) {
	posts, err := ps.postRepository.FindAll()
	if err != nil {
		return nil, err
	}

	response := dto.GetPostsResponseDto{Posts: posts}
	return &response, nil
}

func (ps *PostService) FindById(id string) (*dto.GetPostResponseDto, error) {
	post, err := ps.postRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	response := dto.GetPostResponseDto{Post: *post}
	return &response, nil
}

func (ps *PostService) Update(updateDto dto.UpdatePostRequestDto, postId string) error {
	post, err := ps.postRepository.FindById(postId)
	if err != nil {
		return err
	}
	post.Title = updateDto.Title
	post.Description = updateDto.Description
	post.Status = updateDto.Status
	err = ps.postRepository.Update(post)
	if err != nil {
		return err
	}

	return nil
}
