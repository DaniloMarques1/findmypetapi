package main

import (
	"fmt"
	"log"
	"os"

	"github.com/danilomarques1/findmypetapi/app"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading env variables %v\n", err)
	}
	dbstring := fmt.Sprintf("host=%v dbname=%v password=%v user=%v",
		os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PWD"), os.Getenv("DB_USER"))

	credentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_CONTENT")
	file, err := os.Create(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		log.Fatalf("Error getting google credentials... %v\n", err)
	}

	if _, err := file.WriteString(credentials); err != nil {
		log.Fatalf("Error writing google credentials... %v\n", err)
	}

	var a app.App
	a.Init("database.sql", dbstring)
	a.Listen()
}
