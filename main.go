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

	var a app.App
	a.Init("database.sql", dbstring)
	a.Listen()
}
