package test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"github.com/danilomarques1/findmypetapi/dto"
)

func TestCreateUser(t *testing.T) {
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
}
