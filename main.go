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
	controller, err := vera.GetDeviceRelay(os.Getenv("VERA_DEVICEID"))
	if err != nil {
		log.Println(err)
	}
	//Print out all device names
	for _, device := range controller.SData.Devices {
		log.Println(device.Name)
	}
	controller.Polling()
	for {
		select {
		case <-controller.Updated:
			log.Println("Devices Updated")
			//Print out all device names
			for _, device := range *controller.Switches {
				log.Println(device.Name + "status: " + device.Status)
			}
		}
	}
}
