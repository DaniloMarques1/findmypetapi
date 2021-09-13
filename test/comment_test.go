package test

import (
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/danilomarques1/findmypetapi/dto"
	"github.com/danilomarques1/findmypetapi/model"
	"github.com/danilomarques1/findmypetapi/repository"
	"github.com/danilomarques1/findmypetapi/service"
	"github.com/danilomarques1/findmypetapi/util"
)

type ProducerMock struct {
}

func (p *ProducerMock) Publish(key1, key2 string) error {
	log.Printf("Publishing message...\n")

	return nil
}

func (p *ProducerMock) Setup() error {
	log.Printf("Setting up producer...\n")
	return nil
}

func TestSaveCommentRepository(t *testing.T) {
	cleanTables()
	user := model.User{
		Id:    MOCK_USER_ID,
		Name:  MOCK_USER_NAME,
		Email: MOCK_USER_EMAIL,
	}

	uR := repository.NewUserRepositorySql(App.DB)
	err := uR.Save(&user)
	assertNil(t, err)

	post := model.Post{
		Id:          MOCK_POST1_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post title",
		Description: "Post Description",
		ImageUrl:    "/path/to/image",
		Status:      "missing",
	}

	pRepo := repository.NewPostRepositorySql(App.DB)
	err = pRepo.Save(&post)
	assertNil(t, err)

	comment := model.Comment{
		Id:          MOCK_COMMENT_ID,
		AuthorId:    MOCK_USER_ID,
		PostId:      MOCK_POST1_ID,
		CommentText: "Hey man, i think i saw your dog",
	}

	cRepo := repository.NewCommentRepositorySql(App.DB)
	err = cRepo.Save(&comment)
	assertNil(t, err)
	assertNotEqual(t, "", comment.CreatedAt)
}

func TestSaveCommentService(t *testing.T) {
	cleanTables()
	user := model.User{
		Id:    MOCK_USER_ID,
		Name:  MOCK_USER_NAME,
		Email: MOCK_USER_EMAIL,
	}

	uR := repository.NewUserRepositorySql(App.DB)
	err := uR.Save(&user)
	assertNil(t, err)

	post := model.Post{
		Id:          MOCK_POST1_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post title",
		Description: "Post Description",
		ImageUrl:    "/path/to/image",
		Status:      "missing",
	}

	pRepo := repository.NewPostRepositorySql(App.DB)
	err = pRepo.Save(&post)
	assertNil(t, err)

	requestDto := dto.CreateCommentRequestDto{
		CommentText: "This is a very nice comment",
	}

	cRepo := repository.NewCommentRepositorySql(App.DB)
	producer := ProducerMock{}
	cService := service.NewCommentService(cRepo, &producer)
	response, err := cService.Save(MOCK_USER_ID, MOCK_POST1_ID, requestDto)
	assertNil(t, err)
	assertEqual(t, MOCK_USER_ID, response.Comment.AuthorId)
	assertEqual(t, "This is a very nice comment", response.Comment.CommentText)
}

func TestCreateComment(t *testing.T) {
	cleanTables()
	user := model.User{
		Id:    MOCK_USER_ID,
		Name:  MOCK_USER_NAME,
		Email: MOCK_USER_EMAIL,
	}

	uR := repository.NewUserRepositorySql(App.DB)
	err := uR.Save(&user)
	assertNil(t, err)

	post := model.Post{
		Id:          MOCK_POST1_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post title",
		Description: "Post Description",
		ImageUrl:    "/path/to/image",
		Status:      "missing",
	}

	pRepo := repository.NewPostRepositorySql(App.DB)
	err = pRepo.Save(&post)
	assertNil(t, err)

	token, _, err := util.NewToken(MOCK_USER_ID)
	assertNil(t, err)
	body := `{"comment_text": "This is a very cool comment"}`
	request, err := http.NewRequest(http.MethodPost,
		"/comment/"+MOCK_POST1_ID, strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+token)
	assertNil(t, err)
	response := executeRequest(request)
	assertEqual(t, http.StatusCreated, response.Code)
}

func TestFindCommentsRepository(t *testing.T) {
	cleanTables()
	user := model.User{
		Id:    MOCK_USER_ID,
		Name:  MOCK_USER_NAME,
		Email: MOCK_USER_EMAIL,
	}

	uR := repository.NewUserRepositorySql(App.DB)
	err := uR.Save(&user)
	assertNil(t, err)

	post := model.Post{
		Id:          MOCK_POST1_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post title",
		Description: "Post Description",
		ImageUrl:    "/path/to/image",
		Status:      "missing",
	}

	pRepo := repository.NewPostRepositorySql(App.DB)
	err = pRepo.Save(&post)
	assertNil(t, err)

	comment1 := model.Comment{
		Id:          MOCK_COMMENT_ID,
		AuthorId:    MOCK_USER_ID,
		PostId:      MOCK_POST1_ID,
		CommentText: "New commentary",
	}
	cRepo := repository.NewCommentRepositorySql(App.DB)
	err = cRepo.Save(&comment1)
	assertNil(t, err)

	comments, err := cRepo.FindAll(MOCK_POST1_ID)
	assertNil(t, err)
	assertEqual(t, 1, len(comments))
}

func TestFindCommentsService(t *testing.T) {
	cleanTables()
	user := model.User{
		Id:    MOCK_USER_ID,
		Name:  MOCK_USER_NAME,
		Email: MOCK_USER_EMAIL,
	}

	uR := repository.NewUserRepositorySql(App.DB)
	err := uR.Save(&user)
	assertNil(t, err)

	post := model.Post{
		Id:          MOCK_POST1_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post title",
		Description: "Post Description",
		ImageUrl:    "/path/to/image",
		Status:      "missing",
	}

	pRepo := repository.NewPostRepositorySql(App.DB)
	err = pRepo.Save(&post)
	assertNil(t, err)

	comment1 := model.Comment{
		Id:          MOCK_COMMENT_ID,
		AuthorId:    MOCK_USER_ID,
		PostId:      MOCK_POST1_ID,
		CommentText: "New commentary",
	}
	cRepo := repository.NewCommentRepositorySql(App.DB)
	err = cRepo.Save(&comment1)
	assertNil(t, err)

	cservice := service.NewCommentService(cRepo, &ProducerMock{})
	comments, err := cservice.FindAll(MOCK_POST1_ID)
	assertNil(t, err)
	assertEqual(t, 1, len(comments.Comments))
}

func TestFindComments(t *testing.T) {
	cleanTables()
	user := model.User{
		Id:    MOCK_USER_ID,
		Name:  MOCK_USER_NAME,
		Email: MOCK_USER_EMAIL,
	}

	uR := repository.NewUserRepositorySql(App.DB)
	err := uR.Save(&user)
	assertNil(t, err)

	post := model.Post{
		Id:          MOCK_POST1_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post title",
		Description: "Post Description",
		ImageUrl:    "/path/to/image",
		Status:      "missing",
	}

	pRepo := repository.NewPostRepositorySql(App.DB)
	err = pRepo.Save(&post)
	assertNil(t, err)

	comment1 := model.Comment{
		Id:          MOCK_COMMENT_ID,
		AuthorId:    MOCK_USER_ID,
		PostId:      MOCK_POST1_ID,
		CommentText: "New commentary",
	}
	cRepo := repository.NewCommentRepositorySql(App.DB)
	err = cRepo.Save(&comment1)
	assertNil(t, err)
	request, err := http.NewRequest(http.MethodGet, "/comment/"+MOCK_POST1_ID, nil)
	assertNil(t, err)
	token, _, err := util.NewToken(MOCK_USER_ID)
	assertNil(t, err)
	request.Header.Add("Authorization", "Bearer "+token)
	response := executeRequest(request)
	assertEqual(t, http.StatusOK, response.Code)
}
