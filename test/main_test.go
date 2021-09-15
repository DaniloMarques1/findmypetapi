package test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/danilomarques1/findmypetapi/app"
	"github.com/joho/godotenv"
)

var App app.App

const (
	MOCK_POST1_ID    = "a5886fcf-1de6-462c-8346-d85f72bed0ed"
	MOCK_POST2_ID    = "f25f265b-0c3c-4ecf-a407-675bfa997555"
	MOCK_POST3_ID    = "9e7b5ef7-f28e-4002-bb85-547cca88586b"
	MOCK_USER_ID     = "124e4567-e89b-12d3-a456-426614174000"
	MOCK_USER_NAME   = "Fitz"
	MOCK_USER_EMAIL  = "fitz@gmail.com"
	MOCK_COMMENT_ID  = "dfa3dff3-370d-4618-a04f-ed2dd3b2019b"
	MOCK_USER_ID2    = "f2f4866a-00bc-4db5-9ade-033bb0355b8d"
	MOCK_USER_NAME2  = "Jon"
	MOCK_USER_EMAIL2 = "jon@gmail.com"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading env variables %v\n", err)
	}
	dbstring := fmt.Sprintf("host=%v dbname=%v password=%v user=%v sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PWD"),
		os.Getenv("DB_USER"))

	App.Init("../database.sql", dbstring)

	code := m.Run()
	cleanTables()
	os.Exit(code)
}

func executeRequest(request *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	App.Router.ServeHTTP(rr, request)

	return rr
}

func cleanTables() {
	_, err := App.DB.Exec(`
		truncate table userpet cascade;
		truncate table post cascade;
		truncate table comment cascade;
		`)
	if err != nil {
		log.Fatalf("Error cleaning up the database %v\n", err)
	}
}

func assertEqual(t *testing.T, expect, actual interface{}) {
	if expect != actual {
		t.Fatalf(fmt.Sprintf("\nExpected value: %v\nActual value: %v\n",
			expect, actual))
	}
}

func assertNotEqual(t *testing.T, v1, v2 interface{}) {
	if v1 == v2 {
		t.Fatalf("\nThe values should not be equal\n%v\n%v", v1, v2)
	}
}

func assertNil(t *testing.T, value interface{}) {
	if value != nil {
		t.Fatalf("\nThe given value should be nil\n%v\n", value)
	}
}

func assertNotNil(t *testing.T, value interface{}) {
	if value == nil {
		t.Fatalf("\nThe given value should not be nil\n%v\n", value)
	}
}
