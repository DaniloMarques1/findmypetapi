package test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"github.com/danilomarques1/findmypetapi/dto"
)

func TestCreateUser(t *testing.T) {
	cleanTables()
	body := `{"name": "Fitz", "email": "fitz@gmail.com", "password": "123456", "confirm_password": "123456"}`
	request, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	assertNil(t, err)
	response := executeRequest(request)
	assertEqual(t, http.StatusCreated, response.Code)

	var createResponse dto.CreateUserResponseDto
	err = json.NewDecoder(response.Body).Decode(&createResponse)
	assertNil(t, err)
	assertEqual(t, "Fitz", createResponse.User.Name)
}

func TestErrorCreateUser(t *testing.T) {
	cleanTables()
	// missing name
	body := `{"email": "fitz@gmail.com", "password": "123456", "confirm_password": "123456"}`
	request, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	assertNil(t, err)
	response := executeRequest(request)
	assertEqual(t, http.StatusBadRequest, response.Code)

	// different passwords
	body = `{"name": "Fitz", "email": "fitz@gmail.com", "password": "123456", "confirm_password": "123456890"}`
	request, err = http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	assertNil(t, err)
	response = executeRequest(request)
	assertEqual(t, http.StatusBadRequest, response.Code)

	// large password field
	body = `{"name": "Fitz", "email": "fitz@gmail.com", "password": "39131dkjsakdj12391eqjdhjasdhuh38u29ehdhadhajfhh", "confirm_password": "39131dkjsakdj12391eqjdhjasdhuh38u29ehdhadhajfhh"}`
	request, err = http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	assertNil(t, err)
	response = executeRequest(request)
	assertEqual(t, http.StatusBadRequest, response.Code)
}

func TestCreateUserDuplicate(t *testing.T) {
	cleanTables()
	body := `{"name": "Fitz", "email": "fitz@gmail.com", "password": "123456", "confirm_password": "123456"}`
	request, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	assertNil(t, err)
	response := executeRequest(request)
	assertEqual(t, http.StatusCreated, response.Code)

	// same email, should give an error
	body = `{"name": "Fitz 2", "email": "fitz@gmail.com", "password": "123456", "confirm_password": "123456"}`
	request, err = http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	assertNil(t, err)
	response = executeRequest(request)
	assertEqual(t, http.StatusBadRequest, response.Code)
}

/*
func TestCreateSession(t *testing.T) {
	cleanTables()

	body := `{"name": "Fitz", "email": "fitz@gmail.com", "password": "123456", "confirm_password": "123456"}`
	request, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	assertNil(t, err)
	response := executeRequest(request)
	assertEqual(t, http.StatusCreated, response.Code)

	body = `{"email": "fitz@gmail.com", "password": "123456"}`
	request, err = http.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
	assertNil(t, err)
	response = executeRequest(request)
	assertEqual(t, http.StatusOK, response.Code)
}
*/
