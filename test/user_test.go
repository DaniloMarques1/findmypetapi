package test

import (
	//"log"
	"net/http"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	body := `{"name": "Fitz", "email": "fitz@gmail.com", "password": "123456", "confirm_password": "123456"}`
	request, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(body))
	assertNil(t, err)
	response := executeRequest(request)
	assertEqual(t, http.StatusCreated, response.Code)
}
