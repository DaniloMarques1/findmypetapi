package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type App struct {
	Router    *mux.Router
	validator *validator.Validate
	DB        *sql.DB
}

func (app *App) Init() {
	var err error
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables %v\n", err)
	}

	app.Router = mux.NewRouter()
	app.validator = validator.New()
	app.DB, err = sql.Open("postgres", fmt.Sprintf("host=%v dbname=%v password=%v user=%v sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PWD"), os.Getenv("DB_USER")))
	if err != nil {
		log.Fatalf("Error opening database connection %v\n", err)
	}

	sqlFile, err := os.ReadFile("database.sql")
	if err != nil {
		log.Fatalf("Erro loading databases %v\n", err)
	}
	if _, err := app.DB.Exec(string(sqlFile)); err != nil {
		log.Fatalf("Error creating tables %v\n", err)
	}
}

func (app *App) Listen() {
	port := os.Getenv("SERVER_PORT")
	server := http.Server{
		Handler: app.Router,
		Addr: ":" + port,
	}

	log.Printf("Server starting...")
	log.Fatal(server.ListenAndServe())
}
