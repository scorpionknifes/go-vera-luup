package go-vera-luup

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	conSessionPath = "/info/session/token"
	conRelayPath   = "/relay/relay/relay/device/"
	conDataRequest = "/port_3480/data_request"
	conSData       = "?id=sdata"
	conDevice      = "?id=action&DeviceNum="
	conSwitch      = "urn:upnp-org:serviceId:SwitchPower1"
	conDoorLock    = "urn:micasaverde-com:serviceId:DoorLock1"
)

//GetSessionToken get relay session by using identity
func (con *VeraController) GetSessionToken(identity IdentityJSON) error {
	//Get Url
	url := https + con.ServerRelay + conSessionPath
	//GET Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	//Set Required Headers
	req.Header.Set("MMSAuth", identity.Identity)
	req.Header.Set("MMSAuthSig", identity.IdentitySignature)
	r, err := client.Do(req)
	if err != nil {
		return err
	}
	//Convert response into string as SessionToken
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	con.SessionToken = string(bodyBytes)
	return nil
}

//GetSData Get SData from Hub through Relay Server aka all info
func (con *VeraController) GetSData() error {
	//Get Url
	url := https + con.ServerRelay + conRelayPath + con.DeviceID + conDataRequest + conSData
	//GET Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	//Set Required Headers
	req.Header.Set("MMSSession", con.SessionToken)
	r, err := client.Do(req)
	if err != nil {
		return err
	}
	//Decode SData and add to struct
	sData := SData{}
	err = json.NewDecoder(r.Body).Decode(&sData)
	if err != nil {
		return err
	}
	//Extract SData URL and Session Token for external testing
	//log.Println(url)
	//log.Println(con.SessionToken)
	con.SData = sData

	//Assume DoorLock is a switch

	doorLockID := ""
	switchID := ""
	switches := []Switch{}
	locks := []Lock{}
	for _, category := range sData.Categories {
		if category.Name == "On/Off Switch" {
			switchID = strconv.Itoa(category.ID)
		} else if category.Name == "Doorlock" {
			doorLockID = strconv.Itoa(category.ID)
		}
	}
	if switchID == "" {
		return nil
	}
	for _, device := range sData.Devices {
		if string(device.Category) == switchID {
			deviceID, _ := strconv.Atoi(string(device.ID))
			switchDevice := Switch{Name: device.Name, ID: deviceID, Status: device.Status}
			switches = append(switches, switchDevice)
		}
		if string(device.Category) == doorLockID {
			deviceID, _ := strconv.Atoi(string(device.ID))
			lock := Lock{Name: device.Name, ID: deviceID, Locked: device.Locked}
			locks = append(locks, lock)
		}
	}
	con.Switches = &switches
	con.Locks = &locks
	return nil
}

//Polling loop to CheckStatus using http://wiki.micasaverde.com/index.php/UI_Simple#lu_sdata:_The_polling_loop
func (con *VeraController) Polling() {
	//Loop for polling
	go func() {
		log.Println("Polling")
		poll := Polling{0, 0, 0, con}
		for {
			select {
			case <-con.Kill:
				return
			default:
				err := poll.CheckStatus()
				if err != nil {
					log.Println(err)
					time.Sleep(2 * time.Second)
				}
			}
		}
	}()
}

//pollClient http client without timeout for polling
var pollClient = &http.Client{}

//CheckStatus used to recheck switch status
func (poll *Polling) CheckStatus() error {
	con := poll.VeraController
	//Get Url
	url := https + con.ServerRelay + conRelayPath + con.DeviceID
	params := conDataRequest + conSData + "&loadtime=" + strconv.Itoa(poll.LoadTime) + "&dataversion=" + strconv.Itoa(poll.DataVersion) + "&timeout=60" + "&minimumdelay=60"
	//GET Request
	req, err := http.NewRequest("GET", url+params, nil)
	if err != nil {
		return err
	}
	//Set Required Headers
	req.Header.Set("MMSSession", con.SessionToken)
	r, err := pollClient.Do(req)
	if err != nil {
		return err
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	//Testing Polling data
	//log.Println(string(bodyBytes))
	if string(bodyBytes) == "" {
		time.Sleep(2 * time.Second)
		return nil
	}
	//Decode SData and add to struct
	sData := SData{}
	err = json.Unmarshal(bodyBytes, &sData)
	if err != nil {
		return err
	}

	if sData.Full == 1 {
		//Update all if data is full
		con.SData = sData
	} else {
		//Send con.Updated chan a message when a switch has been updated
		updated := false
		//Update device changes
		for _, d := range sData.Devices {
			//Update switch changes
			for i, e := range *con.Switches {
				id, _ := strconv.Atoi(string(d.ID))
				if id == e.ID {
					if e.Status != d.Status {
						e.Status = d.Status
						sw := *con.Switches
						sw[i] = e
						updated = true
					}
					break
				}
			}
			//Update lock changes
			for i, e := range *con.Locks {
				id, _ := strconv.Atoi(string(d.ID))
				if id == e.ID {
					if e.Locked != d.Locked {
						e.Locked = d.Locked
						lk := *con.Locks
						lk[i] = e
						updated = true
					}
					break
				}
			}
		}
		if updated {
			con.Updated <- true
		}
	}

	//Update polling params read http://wiki.micasaverde.com/index.php/UI_Simple#lu_sdata:_The_polling_loop
	poll.DataVersion = sData.Dataversion
	poll.CurrentMinimumDelay = 2000
	poll.LoadTime = sData.Loadtime
	con.SData = sData
	return nil
}

//Close controller
func (con *VeraController) Close() {
	con.Kill <- true
}

//SwitchPowerStatus change swtich status
func (con *VeraController) SwitchPowerStatus(id int, status int) error {
	url := https + con.ServerRelay + conRelayPath + con.DeviceID
	params := conDataRequest + conDevice + strconv.Itoa(id) + "&serviceId=" + conSwitch + "&action=SetTarget&newTargetValue=" + strconv.Itoa(status)
	err := con.CallURL(url + params)
	if err != nil {
		return err
	}
	return nil
}

//DoorLockStatus change lock status
func (con *VeraController) DoorLockStatus(id int, status int) error {
	url := https + con.ServerRelay + conRelayPath + con.DeviceID
	params := conDataRequest + conDevice + strconv.Itoa(id) + "&serviceId=" + conDoorLock + "&action=SetTarget&newTargetValue=" + strconv.Itoa(status)
	err := con.CallURL(url + params)
	if err != nil {
		return err
	}
	return nil
}

//CallURL using GET
func (con *VeraController) CallURL(url string) error {
	//GET Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	//Set Required Headers
	req.Header.Set("MMSSession", con.SessionToken)
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
