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
	vera := New(os.Getenv("VERA_USERNAME"), os.Getenv("VERA_PASSWORD"))

	//DeviceID = SN number on Vera controller
	deviceInfo, err := vera.GetDeviceInfo(os.Getenv("VERA_DEVICEID"))
	if err != nil {
		log.Println(err)
	}
	log.Println(deviceInfo)
}
