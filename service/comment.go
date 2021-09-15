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
	// TODO produce rabbit mq message
	commentId := uuid.NewString()
	comment := model.Comment{
		Id:          commentId,
		AuthorId:    userId,
		PostId:      postId,
		CommentText: request.CommentText,
	}
	err := cs.commentRepository.Save(&comment)
	if err != nil {
		return nil, err
	}
	response := dto.CreateCommentResponseDto{
		Comment: comment,
	}

	go func() {
		msg, err := cs.commentRepository.GetCommentNotificationMessage(postId, commentId)
		if err == nil && len(msg) != 0 {
			err = cs.producer.Publish(msg)
			if err != nil {
				// TODO how to handle errors
				log.Printf("Error publishing message %v\n", err)
			}
		}
	}()

	return &response, nil
}

func (cs *CommentService) FindAll(postId string) (*dto.GetCommentsResponseDto, error) {
	comments, err := cs.commentRepository.FindAll(postId)
	if err != nil {
		return nil, err
	}
	response := dto.GetCommentsResponseDto{Comments: comments}

	return &response, nil
}
