package test

import (
	"encoding/json"
	"github.com/danilomarques1/findmypetapi/dto"
	"net/http"
	"strings"
	"testing"
)

// testing controllers

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

	var sessionResponse dto.SessionResponseDto
	err = json.NewDecoder(response.Body).Decode(&sessionResponse)

	assertNil(t, err)
	assertNotNil(t, sessionResponse.User)
	assertNotEqual(t, "", sessionResponse.Token)
	assertNotEqual(t, "", sessionResponse.RefreshToken)
	assertNotEqual(t, "", sessionResponse.User.Name)
}

func TestCreateSessionError(t *testing.T) {
	cleanTables()

	body := `{"name": "Fitz", "email": "fitz@gmail.com", "password": "123456", "confirm_password": "123456"}`
	request, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	assertNil(t, err)
	response := executeRequest(request)
	assertEqual(t, http.StatusCreated, response.Code)

	// not sending required field email
	body = `{password": "123456"}`
	request, err = http.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
	assertNil(t, err)
	response = executeRequest(request)
	assertEqual(t, http.StatusBadRequest, response.Code)
}

func TestRefreshSession(t *testing.T) {
	cleanTables()

	body := `{"name": "Fitz", "email": "fitz@gmail.com", "password": "123456", "confirm_password": "123456"}`
	request, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	assertNil(t, err)
	response := executeRequest(request)
	assertEqual(t, http.StatusCreated, response.Code)

	body = `{"email": "fitz@gmail.com", "password":"123456"}`
	request, err = http.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
	assertNil(t, err)
	response = executeRequest(request)
	assertEqual(t, http.StatusOK, response.Code)

	var responseDto dto.SessionResponseDto
	err = json.NewDecoder(response.Body).Decode(&responseDto)
	assertNil(t, err)

	request, err = http.NewRequest(http.MethodPut, "/session/refresh", nil)
	assertNil(t, err)
	request.Header.Add("refresh_token", responseDto.RefreshToken)
	response = executeRequest(request)
	assertEqual(t, http.StatusOK, response.Code)
}

func TestRefreshSessionError(t *testing.T) {
	cleanTables()

	body := `{"name": "Fitz", "email": "fitz@gmail.com", "password": "123456", "confirm_password": "123456"}`
	request, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	assertNil(t, err)
	response := executeRequest(request)
	assertEqual(t, http.StatusCreated, response.Code)

	body = `{"email": "fitz@gmail.com", "password":"123456"}`
	request, err = http.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
	assertNil(t, err)
	response = executeRequest(request)
	assertEqual(t, http.StatusOK, response.Code)

	var responseDto dto.SessionResponseDto
	err = json.NewDecoder(response.Body).Decode(&responseDto)
	assertNil(t, err)

	// using token instead of refresh token
	request, err = http.NewRequest(http.MethodPut, "/session/refresh", nil)
	assertNil(t, err)
	request.Header.Add("refresh_token", responseDto.Token)
	response = executeRequest(request)
	assertEqual(t, http.StatusUnauthorized, response.Code)

	// sending empty string
	request, err = http.NewRequest(http.MethodPut, "/session/refresh", nil)
	assertNil(t, err)
	request.Header.Add("refresh_token", "")
	response = executeRequest(request)
	assertEqual(t, http.StatusUnauthorized, response.Code)

	// not sending header at all
	request, err = http.NewRequest(http.MethodPut, "/session/refresh", nil)
	assertNil(t, err)
	response = executeRequest(request)
	assertEqual(t, http.StatusUnauthorized, response.Code)
}

func TestUpdateUser(t *testing.T) {
	cleanTables()

	body := `{"name": "Fitz", "email": "fitz@gmail.com", "password": "123456", "confirm_password": "123456"}`
	request, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	assertNil(t, err)
	response := executeRequest(request)
	assertEqual(t, http.StatusCreated, response.Code)

	body = `{"email": "fitz@gmail.com", "password":"123456"}`
	request, err = http.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
	assertNil(t, err)
	response = executeRequest(request)
	assertEqual(t, http.StatusOK, response.Code)
	var sessionResponse dto.SessionResponseDto
	err = json.NewDecoder(response.Body).Decode(&sessionResponse)
	assertNil(t, err)

	body = `{"name": "Fitz New", "old_password": "123456", "new_password": "newpassword",
		"confirm_password": "newpassword"}`
	request, err = http.NewRequest(http.MethodPut, "/user", strings.NewReader(body))
	request.Header.Add("Authorization", "Bearer "+sessionResponse.Token)
	assertNil(t, err)

	response = executeRequest(request)
	assertEqual(t, http.StatusNoContent, response.Code)

	body = `{"email": "fitz@gmail.com", "password":"newpassword"}`
	request, err = http.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
	assertNil(t, err)
	response = executeRequest(request)
	assertEqual(t, http.StatusOK, response.Code)
}

func TestUpdateUserError(t *testing.T) {
	cleanTables()

	body := `{"name": "Fitz", "email": "fitz@gmail.com", "password": "123456", "confirm_password": "123456"}`
	request, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	assertNil(t, err)
	response := executeRequest(request)
	assertEqual(t, http.StatusCreated, response.Code)

	body = `{"email": "fitz@gmail.com", "password":"123456"}`
	request, err = http.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
	assertNil(t, err)
	response = executeRequest(request)
	assertEqual(t, http.StatusOK, response.Code)
	var sessionResponse dto.SessionResponseDto
	err = json.NewDecoder(response.Body).Decode(&sessionResponse)
	assertNil(t, err)

	// wrong old password, 400
	body = `{"name": "Fitz", "old_password": "wrongpassword", 
		"new_password": "validnewpassword", "confirm_password": "validnewpassword"}`
	request, err = http.NewRequest(http.MethodPut, "/user", strings.NewReader(body))
	request.Header.Add("Authorization", "Bearer "+sessionResponse.Token)
	assertNil(t, err)
	response = executeRequest(request)
	assertEqual(t, http.StatusBadRequest, response.Code)

	// missing required field, 400
	body = `{"old_password": "wrongpassword", 
		"new_password": "validnewpassword", "confirm_password": "validnewpassword"}`
	request, err = http.NewRequest(http.MethodPut, "/user", strings.NewReader(body))
	request.Header.Add("Authorization", "Bearer "+sessionResponse.Token)
	assertNil(t, err)
	response = executeRequest(request)
	assertEqual(t, http.StatusBadRequest, response.Code)

	// missing token, 401
	body = `{ "name": "Fitz", "old_password": "wrongpassword", 
		"new_password": "validnewpassword", "confirm_password": "validnewpassword"}`
	request, err = http.NewRequest(http.MethodPut, "/user", strings.NewReader(body))
	assertNil(t, err)
	response = executeRequest(request)
	assertEqual(t, http.StatusUnauthorized, response.Code)
}

// testing repository

// testing services
