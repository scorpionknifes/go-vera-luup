package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	//Load .env config - username and password
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//Example Create new object e.g vera = New(username, password)
	New(os.Getenv("VERAUSERNAME"), os.Getenv("VERAPASSWORD"))
}
