package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router    *mux.Router
	validator *validator.Validate
	DB        *sql.DB
}

func (app *App) Init(sqlFileName, dbstring string) {
	var err error

	app.Router = mux.NewRouter()
	app.validator = validator.New()
	app.DB, err = sql.Open("postgres", dbstring)
	if err != nil {
		log.Fatalf("Error opening database connection %v\n", err)
	}

	sqlFile, err := os.ReadFile(sqlFileName)
	if err != nil {
		log.Fatalf("Error loading databases %v\n", err)
	}
	if _, err := app.DB.Exec(string(sqlFile)); err != nil {
		log.Fatalf("Error creating tables %v\n", err)
	}
}

func (app *App) Listen() {
	port := os.Getenv("SERVER_PORT")
	server := http.Server{
		Handler: app.Router,
		Addr:    ":" + port,
	}

	log.Printf("Server starting...")
	log.Fatal(server.ListenAndServe())
}
