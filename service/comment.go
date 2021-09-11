package service

import (
	"github.com/danilomarques1/findmypetapi/dto"
	"github.com/danilomarques1/findmypetapi/model"
	"github.com/google/uuid"
)

type CommentService struct {
	commentRepository model.CommentRepository
}

func NewCommentService(commentRepository model.CommentRepository) *CommentService {
	return &CommentService{
		commentRepository: commentRepository,
	}
}

func (cs *CommentService) Save(userId, postId string,
	request dto.CreateCommentRequestDto) (*dto.CreateCommentResponseDto, error) {
	// TODO produce rabbit mq message
	id := uuid.NewString()
	comment := model.Comment{
		Id:          id,
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
