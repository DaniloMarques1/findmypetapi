package test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/danilomarques1/findmypetapi/dto"
	"github.com/danilomarques1/findmypetapi/model"
	"github.com/danilomarques1/findmypetapi/repository"
	"github.com/danilomarques1/findmypetapi/service"
	"github.com/danilomarques1/findmypetapi/util"
)

const MOCK_USER_ID = "124e4567-e89b-12d3-a456-426614174000"
const MOCK_USER_NAME = "Fitz"
const MOCK_USER_EMAIL = "fitz@gmail.com"

// testing the post repository
func TestSavePostRepository(t *testing.T) {
	cleanTables()
	userRepository := repository.NewUserRepositorySql(App.DB)
	user := model.User{
		Id:    MOCK_USER_ID,
		Name:  MOCK_USER_NAME,
		Email: MOCK_USER_EMAIL,
	}
	err := userRepository.Save(&user)
	assertNil(t, err)

	postRepository := repository.NewPostRepositorySql(App.DB)
	post := model.Post{
		Id:          "123e4567-e89b-12d3-a456-426614174000",
		AuthorId:    "124e4567-e89b-12d3-a456-426614174000",
		Title:       "Post title",
		Description: "description",
		ImageUrl:    "some/path/to/file",
		Status:      "missing",
		CreatedAt:   time.Now(),
	}

	err = postRepository.Save(&post)
	assertNil(t, err)
}

func TestCreatePostService(t *testing.T) {
	cleanTables()
	userRepository := repository.NewUserRepositorySql(App.DB)
	user := model.User{
		Id:    MOCK_USER_ID,
		Name:  MOCK_USER_NAME,
		Email: MOCK_USER_EMAIL,
	}
	err := userRepository.Save(&user)
	assertNil(t, err)

	postRepository := repository.NewPostRepositorySql(App.DB)
	postService := service.NewPostService(postRepository)
	postDto := dto.CreatePostRequestDto{
		Title:       "Post title",
		Description: "Post description",
	}

	response, err := postService.CreatePost(postDto, MOCK_USER_ID)
	assertNil(t, err)
	assertEqual(t, "Post title", response.Post.Title)
	assertEqual(t, MOCK_USER_ID, response.Post.AuthorId)
}

func TestCreatePost(t *testing.T) {
	cleanTables()
	userRepository := repository.NewUserRepositorySql(App.DB)
	user := model.User{
		Id:    MOCK_USER_ID,
		Name:  MOCK_USER_NAME,
		Email: MOCK_USER_EMAIL,
	}
	err := userRepository.Save(&user)
	assertNil(t, err)

	body := `
		{
			"title": "Post title",
			"description": "Post description"
		}
	`
	token, _, err := util.NewToken(MOCK_USER_ID)
	assertNil(t, err)
	assertNotEqual(t, "", token)

	request, err := http.NewRequest(http.MethodPost, "/post", strings.NewReader(body))
	assertNil(t, err)
	request.Header.Add("Authorization", "Bearer "+token)
	response := executeRequest(request)
	assertEqual(t, http.StatusCreated, response.Code)

	var dto dto.CreatePostResponseDto
	err = json.NewDecoder(response.Body).Decode(&dto)
	assertNil(t, err)
	assertEqual(t, MOCK_USER_ID, dto.Post.AuthorId)
	assertEqual(t, "missing", dto.Post.Status)
}
