package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	vera := Vera{
		Username: os.Getenv("VERAUSERNAME"),
		Password: os.Getenv("VERAPASSWORD"),
	}

	vera.GetLoginToken()
}
