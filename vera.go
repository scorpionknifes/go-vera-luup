package govera

import (
	"log"
	"time"
)

//New Create new Vera object
func New(username string, password string) Vera {
	//Initialise Object
	vera := Vera{
		Username:    username,
		Password:    password,
		Controllers: &[]VeraController{},
	}
	// Setup Identity, SessionToken
	err := vera.Renew()
	if err != nil {
		log.Panic(err)
	}
	//Gets all devices linked to account using SessionToken
	err = vera.GetAllDevices()
	if err != nil {
		log.Panic(err)
	}

	//Loop 23 hrs to keep renewing Tokens
	ticker := time.NewTicker(23 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				ticker = time.NewTicker(23 * time.Hour)
				err = vera.Renew()
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	return vera
}

//Renew Used to renew identity
func (vera *Vera) Renew() error {
	//Gets Identity using username and password
	err := vera.GetIdentityToken()
	if err != nil {
		return err
	}
	//Gets SessionToken using Identity
	err = vera.GetSessionToken()
	if err != nil {
		return err
	}

	//Renew all controllers
	log.Println("Renewed")
	for _, controller := range *vera.Controllers {
		err = controller.GetSessionToken()
		if err != nil {
			vera.RemoveDevice(controller.DeviceID)
		}
	}

	return nil
}
