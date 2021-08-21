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

var a app.App

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading env variables %v\n", err)
	}
	dbstring := fmt.Sprintf("host=%v dbname=%v password=%v user=%v sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PWD"), os.Getenv("DB_USER"))

	a.Init("../database.sql", dbstring)

	code := m.Run()
	cleanTables()
	os.Exit(code)
}

func executeRequest(request *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, request)

	return rr
}

func cleanTables() {
	_, err := a.DB.Exec("truncate table userpet cascade; truncate table userpost cascade; truncate table post cascade;")
	if err != nil {
		log.Fatalf("Error cleaning up the database %v\n", err)
	}
}

func assertEqual(t *testing.T, expect, actual interface{}) {
	if expect != actual {
		t.Fatalf(fmt.Sprintf("\nExpected value: %v\nActual value: %v\n", expect, actual))
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
