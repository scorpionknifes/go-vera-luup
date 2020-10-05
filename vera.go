// Package vera to remotely access Vera™ home controller UI7 and call Luup
// For more info about Luup http://wiki.micasaverde.com/index.php/Luup_Requests
//
// Example
//
//	//Example Create new object e.g vera = New(username, password)
//	user := vera.New("username", "password")
//
//	//DeviceID = SN number on Vera controller
//	controller, err := user.GetDeviceRelay("12345670")
//	if err != nil {
//		log.Println(err)
//	}
//
//	//Close controller by
//	controller.Close()
//
//	//Change Switch ID: 5 to Status: 1 aka Turn on Switch 5
//	controller.SwitchPowerStatus(5, 1)
//
//	//Lock door
//	controller.DoorLockStatus(LockID, 1) // 1 = lock, 0 = unlock
//
//	//Check Status using go channels
//	for {
//		select {
//		case <-controller.Updated:
//			log.Println("Devices Updated")
//			//Print out all device names
//			for _, device := range *controller.Switches {
//				log.Println("Device: " + device.Name + " ID: " + strconv.Itoa(device.ID) + " status: " + device.Status)
//			}
//
//			for _, lock := range *controller.Locks {
//				log.Println("Lock: " + lock.Name + " ID: " + strconv.Itoa(lock.ID) + " status: " + lock.Locked)
//			}
//		}
//	}
//
package vera

import (
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	// https url.
	https = "https://"

	// remote url mios.
	remoteURL    = "us-autha11.mios.com"
	passwordSeed = "oZ7QE6LcLJp6fiWzdqZc" //nolint:gosec

	// url paths.
	loginPath      = "/autha/auth/username/"
	sessionPath    = "/info/session/token"
	devicePath     = "/device/device/device/"
	accountPath    = "/account/account/account/"
	devicesPath    = "/devices"
	conSessionPath = "/info/session/token"
	conRelayPath   = "/relay/relay/relay/device/"
	conDataRequest = "/port_3480/data_request"

	// url params.
	conSData  = "?id=sdata"
	conDevice = "?id=action&DeviceNum="

	// serviceId.
	conSwitch   = "urn:upnp-org:serviceId:SwitchPower1"
	conDoorLock = "urn:micasaverde-com:serviceId:DoorLock1"

	// time.
	renewTimer    = 23 * time.Hour
	clientTimeout = 10 * time.Second
	retryTimer    = 5 * time.Second
)

// client with 10 sec timeout.
func client() *http.Client {
	return &http.Client{Timeout: clientTimeout}
}

// New Create new Vera object that identify using username & password
// Login using Vera™ home controller UI7 account.
func New(username string, password string) Vera {
	// Initialize Vera Class Object
	vera := Vera{
		Username:    username,
		Password:    password,
		Controllers: &[]Controller{},
		m:           &sync.Mutex{},
	}
	// Setup Identity, SessionToken
	err := vera.Renew()
	if err != nil {
		log.Panic(err)
	}
	// Gets all devices linked to account using SessionToken
	err = vera.GetAllDevices()
	if err != nil {
		log.Panic(err)
	}

	// Identity expires in 24 hr after login
	// Loop 23 hrs to keep renewing Tokens
	ticker := time.NewTicker(renewTimer)

	go func() {
		for {
			<-ticker.C
			ticker = time.NewTicker(renewTimer)

			err = vera.Renew()
			if err != nil {
				log.Println(err)
			}
		}
	}()

	return vera
}

// Renew Used to renew identity of Vera struct
// Call this Renew() to manually renew Vera Identity.
func (vera *Vera) Renew() error {
	vera.m.Lock()

	// Renew Identity using username and password
	err := vera.GetIdentityToken()
	if err != nil {
		return err
	}

	log.Println("GetIdentity")

	// Renew SessionToken using Identity
	err = vera.GetSessionToken()
	if err != nil {
		return err
	}

	// Renew all controllers
	log.Println("Renewed")

	for _, controller := range *vera.Controllers {
		// Remove controller if error when renewing
		err = controller.Renew(*vera)
		if err != nil {
			log.Println("Removed " + controller.DeviceID)

			err = vera.removeDevice(controller.DeviceID)
			if err != nil {
				return err
			}
		}
	}

	log.Println("Complete")
	vera.m.Unlock()

	return nil
}
