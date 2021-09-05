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

const (
	MOCK_POST1_ID   = "a5886fcf-1de6-462c-8346-d85f72bed0ed"
	MOCK_POST2_ID   = "f25f265b-0c3c-4ecf-a407-675bfa997555"
	MOCK_POST3_ID   = "9e7b5ef7-f28e-4002-bb85-547cca88586b"
	MOCK_USER_ID    = "124e4567-e89b-12d3-a456-426614174000"
	MOCK_USER_NAME  = "Fitz"
	MOCK_USER_EMAIL = "fitz@gmail.com"
)

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
		AuthorId:    MOCK_USER_ID,
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

func TestFindAllRepository(t *testing.T) {
	cleanTables()
	userRepository := repository.NewUserRepositorySql(App.DB)
	user := model.User{Id: MOCK_USER_ID, Name: "Fitz", Email: "fitz@gmail.com"}
	err := userRepository.Save(&user)
	assertNil(t, err)

	postRepository := repository.NewPostRepositorySql(App.DB)
	post1 := model.Post{
		Id:          MOCK_POST1_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post 1",
		Description: "Desc 1",
		ImageUrl:    "/path/to/image",
	}
	err = postRepository.Save(&post1)
	assertNil(t, err)

	post2 := model.Post{
		Id:          MOCK_POST2_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post 2",
		Description: "Desc 2",
		ImageUrl:    "/path/to/image",
	}
	err = postRepository.Save(&post2)
	assertNil(t, err)

	post3 := model.Post{
		Id:          MOCK_POST3_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post 3",
		Description: "Desc 3",
		ImageUrl:    "/path/to/image",
	}
	err = postRepository.Save(&post3)
	assertNil(t, err)

	posts, err := postRepository.FindAll()

	assertNil(t, err)
	assertNotNil(t, posts)
	assertEqual(t, 3, len(posts))
	assertEqual(t, "Post 1", posts[0].Title)
}

func TestGetAllPostsService(t *testing.T) {
	cleanTables()
	userRepository := repository.NewUserRepositorySql(App.DB)
	user := model.User{Id: MOCK_USER_ID, Name: "Fitz", Email: "fitz@gmail.com"}
	err := userRepository.Save(&user)
	assertNil(t, err)

	postRepository := repository.NewPostRepositorySql(App.DB)
	post1 := model.Post{
		Id:          MOCK_POST1_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post 1",
		Description: "Desc 1",
		ImageUrl:    "/path/to/image",
	}
	err = postRepository.Save(&post1)
	assertNil(t, err)

	post2 := model.Post{
		Id:          MOCK_POST2_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post 2",
		Description: "Desc 2",
		ImageUrl:    "/path/to/image",
	}
	err = postRepository.Save(&post2)
	assertNil(t, err)

	post3 := model.Post{
		Id:          MOCK_POST3_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post 3",
		Description: "Desc 3",
		ImageUrl:    "/path/to/image",
	}
	err = postRepository.Save(&post3)
	assertNil(t, err)

	postService := service.NewPostService(postRepository)
	response, err := postService.GetAll()
	assertNil(t, err)
	assertEqual(t, 3, len(response.Posts))
}

func TestGetAll(t *testing.T) {
	cleanTables()
	userRepository := repository.NewUserRepositorySql(App.DB)
	user := model.User{Id: MOCK_USER_ID, Name: "Fitz", Email: "fitz@gmail.com"}
	err := userRepository.Save(&user)
	assertNil(t, err)

	postRepository := repository.NewPostRepositorySql(App.DB)
	post1 := model.Post{
		Id:          MOCK_POST1_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post 1",
		Description: "Desc 1",
		ImageUrl:    "/path/to/image",
	}
	err = postRepository.Save(&post1)
	assertNil(t, err)

	post2 := model.Post{
		Id:          MOCK_POST2_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post 2",
		Description: "Desc 2",
		ImageUrl:    "/path/to/image",
	}
	err = postRepository.Save(&post2)
	assertNil(t, err)

	post3 := model.Post{
		Id:          MOCK_POST3_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post 3",
		Description: "Desc 3",
		ImageUrl:    "/path/to/image",
	}
	err = postRepository.Save(&post3)
	assertNil(t, err)

	token, _, err := util.NewToken(MOCK_USER_ID)
	assertNil(t, err)

	request, err := http.NewRequest(http.MethodGet, "/post", nil)
	assertNil(t, err)
	request.Header.Add("Authorization", "Bearer "+token)
	response := executeRequest(request)
	assertEqual(t, http.StatusOK, response.Code)

	var posts dto.GetPostResponseDto
	err = json.NewDecoder(response.Body).Decode(&posts)
	assertNil(t, err)
	assertEqual(t, 3, len(posts.Posts))
}

func TestFindById(t *testing.T) {
	cleanTables()
	user := model.User{Id: MOCK_USER_ID, Name: "Fitz", Email: "fitz@gmail.com"}
	postToBeCreated := model.Post{
		Id:          MOCK_POST1_ID,
		AuthorId:    MOCK_USER_ID,
		Title:       "Post title",
		Description: "Post Description",
		ImageUrl:    "/path/to/image",
		Status:      "missing",
	}

	pRepo := repository.NewPostRepositorySql(App.DB)
	uRepo := repository.NewUserRepositorySql(App.DB)

	err := uRepo.Save(&user)
	assertNil(t, err)

	err = pRepo.Save(&postToBeCreated)
	assertNil(t, err)

	foundP, err := pRepo.FindById(MOCK_POST1_ID)
	assertNil(t, err)
	assertNotNil(t, foundP)
	assertEqual(t, foundP.Title, "Post title")
	assertEqual(t, foundP.AuthorId, MOCK_USER_ID)
}
