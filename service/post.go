package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/danilomarques1/findmypetapi/dto"
	"github.com/danilomarques1/findmypetapi/lib"
	"github.com/danilomarques1/findmypetapi/model"
	"github.com/danilomarques1/findmypetapi/util"
	"github.com/google/uuid"
)

type PostService struct {
	postRepository model.PostRepository
	producer       lib.Producer
}

func NewPostService(postRepository model.PostRepository, producer lib.Producer) *PostService {
	return &PostService{
		postRepository: postRepository,
		producer:       producer,
	}
}

func (ps *PostService) CreatePost(postDto dto.CreatePostRequestDto,
	userId string) (*dto.CreatePostResponseDto, error) {
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

func (ps *PostService) Update(updateDto dto.UpdatePostRequestDto, authorId, postId string) error {
	post, err := ps.postRepository.FindPostByAuthor(authorId, postId)
	if err != nil {
		log.Printf("No posts were found for this user")
		return util.NewApiError("Post not found", http.StatusNotFound)
	}

	log.Printf("%v\n", updateDto)
	if post.Status != updateDto.Status && updateDto.Status == "found" {
		msg := dto.StatusChangeNotification{PostId: postId}
		mBytes, err := json.Marshal(&msg)
		log.Printf("Marshal message %v\n", err)
		if err == nil {
			log.Printf("Message %v\n", string(mBytes))
			err = ps.producer.Publish(mBytes, lib.STATUS_CHANGE_QUEUE)
			if err != nil {
				log.Printf("Error publishing %v\n", err)
			}
		}
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
