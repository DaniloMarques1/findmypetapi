package service

import (
	"log"

	"github.com/danilomarques1/findmypetapi/dto"
	"github.com/danilomarques1/findmypetapi/lib"
	"github.com/danilomarques1/findmypetapi/model"
	"github.com/google/uuid"
)

type CommentService struct {
	commentRepository model.CommentRepository
	producer          lib.Producer
}

func NewCommentService(commentRepository model.CommentRepository,
	producer lib.Producer) *CommentService {
	return &CommentService{
		commentRepository: commentRepository,
		producer:          producer,
	}
}

func (cs *CommentService) Save(userId, postId string,
	request dto.CreateCommentRequestDto) (*dto.CreateCommentResponseDto, error) {
	commentId := uuid.NewString()
	post := &model.Post{Id: postId}
	author := &model.User{Id: userId}
	comment := model.Comment{
		Id:          commentId,
		Author:      author,
		Post:        post,
		CommentText: request.CommentText,
	}
	err := cs.commentRepository.Save(&comment)
	if err != nil {
		return nil, err
	}
	response := dto.CreateCommentResponseDto{
		Comment: dto.CommentDto{
			Id:          comment.Id,
			CreatedAt:   comment.CreatedAt,
			CommentText: comment.CommentText,
		},
	}

	msg, err := cs.commentRepository.GetCommentNotificationMessage(postId, commentId)
	if err == nil && len(msg) != 0 {
		err = cs.producer.Publish(msg, lib.COMMENT_QUEUE)
		if err != nil {
			// TODO how to handle errors
			log.Printf("Error publishing message %v\n", err)
		}
	}

	return &response, nil
}

func (cs *CommentService) FindAll(postId string) (*dto.GetCommentsResponseDto, error) {
	comments, err := cs.commentRepository.FindAll(postId)
	if err != nil {
		return nil, err
	}

	response := dto.GetCommentsResponseDto{
		Comments: comments,
	}

	return &response, nil
}
