package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	vera "github.com/scorpionknifes/go-vera-luup"
)

func main() {
	//Load .env config - username and password
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//Example Create new object e.g vera = New(username, password)
	user := vera.New(os.Getenv("VERA_USERNAME"), os.Getenv("VERA_PASSWORD"))

	//DeviceID = SN number on Vera controller
	controller, err := user.GetDeviceRelay(os.Getenv("VERA_DEVICEID"))
	if err != nil {
		log.Println(err)
	}

	//Close controller by
	controller.Close()

	//Change Switch ID: 5 to Status: 1 aka Turn on Switch 5
	controller.SwitchPowerStatus(5, 1)

	//Lock door
	lockID, _ := strconv.Atoi(os.Getenv("VERA_LOCKID"))
	controller.DoorLockStatus(lockID, 1) // 1 = lock, 0 = unlock

	//Check Status using go channels
	for {
		select {
		case <-controller.Updated:
			log.Println("Devices Updated")
			//Print out all device names
			for _, device := range *controller.Switches {
				log.Println("Device: " + device.Name + " ID: " + strconv.Itoa(device.ID) + " status: " + device.Status)
			}

			for _, lock := range *controller.Locks {
				log.Println("Lock: " + lock.Name + " ID: " + strconv.Itoa(lock.ID) + " status: " + lock.Locked)
			}
		}
	}
}
