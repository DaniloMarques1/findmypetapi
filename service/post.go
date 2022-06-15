package service

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
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

const (
	BUCKET_NAME = "findmypetbucket"
	BUCKET_URL  = "https://storage.googleapis.com/findmypetbucket/"
)

func (ps *PostService) CreatePost(postDto dto.CreatePostRequestDto,
	userId string) (*dto.CreatePostResponseDto, error) {
	path, err := ps.uploadFile(postDto.Filename, postDto.File)
	if err != nil {
		log.Printf("Error uploading file %v\n", err)
		return nil, err
	}

	postId := uuid.NewString()
	author := &model.User{Id: userId}
	post := model.Post{
		Id:          postId,
		Author:      author,
		Title:       postDto.Title,
		Description: postDto.Description,
		ImageUrl:    path,
	}

	err = ps.postRepository.Save(&post)
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

func (ps *PostService) FindPostsByAuthor(authorId string) (*dto.GetPostsResponseDto, error) {
	posts, err := ps.postRepository.FindPostsByAuthor(authorId)
	if err != nil {
		return nil, err
	}

	response := dto.GetPostsResponseDto{
		Posts: posts,
	}

	return &response, nil
}

func (ps *PostService) uploadFile(fileName string, file multipart.File) (string, error) {
	return "/path/to/file", nil // TODO remove
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	bucket := client.Bucket(BUCKET_NAME)
	object := bucket.Object(fileName)

	wc := object.NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	return BUCKET_URL + fileName, nil
}
