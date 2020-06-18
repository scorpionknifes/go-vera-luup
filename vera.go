//Package vera to remotely accesss Vera™ home controller UI7 and call Luup
//For more info about Luup http://wiki.micasaverde.com/index.php/Luup_Requests
//
//Example
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
	"time"
)

const (
	https          = "https://"
	remoteURL      = "us-autha11.mios.com"
	loginPath      = "/autha/auth/username/"
	passwordSeed   = "oZ7QE6LcLJp6fiWzdqZc"
	sessionPath    = "/info/session/token"
	devicePath     = "/device/device/device/"
	accountPath    = "/account/account/account/"
	devicesPath    = "/devices"
	conSessionPath = "/info/session/token"
	conRelayPath   = "/relay/relay/relay/device/"
	conDataRequest = "/port_3480/data_request"
	conSData       = "?id=sdata"
	conDevice      = "?id=action&DeviceNum="
	conSwitch      = "urn:upnp-org:serviceId:SwitchPower1"
	conDoorLock    = "urn:micasaverde-com:serviceId:DoorLock1"
)

var (
	client = &http.Client{Timeout: 10 * time.Second}
	//pollClient http client without timeout for polling
	pollClient = &http.Client{}
)

//New Create new Vera object
func New(username string, password string) Vera {
	//Initialise Object
	vera := Vera{
		Username:    username,
		Password:    password,
		Controllers: &[]Controller{},
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

//Renew Used to renew identity of Vera struct
func (vera *Vera) Renew() error {
	//Renew Identity using username and password
	err := vera.GetIdentityToken()
	if err != nil {
		return err
	}
	//Renew SessionToken using Identity
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
