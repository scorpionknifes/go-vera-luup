package vera

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Renew renews controller by getting sessions
// Kills polling and restart polling.
func (con *Controller) Renew(vera Vera) error {
	con.Kill <- true
	con.Vera = &vera

	err := con.GetSessionToken(vera)
	if err != nil {
		return err
	}

	con.Polling()

	return nil
}

// GetSessionToken get relay session by using identity
// Call GetSessionToken() to manually renew session token.
func (con *Controller) GetSessionToken(vera Vera) error {
	identity := vera.Identity
	url := https + con.ServerRelay + conSessionPath

	// GET Request
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	// Set Required Headers
	req.Header.Set("MMSAuth", identity.Identity)
	req.Header.Set("MMSAuthSig", identity.IdentitySignature)

	r, err := client().Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// Convert response into string as SessionToken
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	con.SessionToken = string(bodyBytes)

	return nil
}

// GetSData Get SData from Hub through Relay Server aka all info
// Info is stored back inside vera. This should be only called to get new SData.
func (con *Controller) GetSData() error { //nolint:funlen
	// Get Url
	url := https + con.ServerRelay + conRelayPath + con.DeviceID + conDataRequest + conSData
	// GET Request
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	// Set Required Headers
	req.Header.Set("MMSSession", con.SessionToken)

	r, err := client().Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// Decode SData and add to struct
	sData := SData{}

	err = json.NewDecoder(r.Body).Decode(&sData)
	if err != nil {
		return err
	}
	// Extract SData URL and Session Token for external testing
	// log.Println(url)
	// log.Println(con.SessionToken)
	con.SData = sData

	// Assume DoorLock is a switch

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

// Polling using http://wiki.micasaverde.com/index.php/UI_Simple#lu_sdata:_The_polling_loop
// Polling is achieved by calling GET request and server side timeout until a change is detected.
func (con *Controller) Polling() {
	// Loop for polling
	go func() {
		// log.Println("Polling")
		poll := Polling{0, 0, 0, con}

		for {
			select {
			case <-con.Kill:
				return
			default:
				err := poll.checkStatus()
				if err != nil {
					log.Println("Retry in 5 sec")
					time.Sleep(retryTimer)
				}
			}
		}
	}()
}

// checkStatus calls GET request for polling, go channels used to signal when data has been changed.
func (poll *Polling) checkStatus() error { //nolint:funlen,gocognit
	con := poll.Controller
	// Get Url
	url := https + con.ServerRelay + conRelayPath + con.DeviceID
	dataParams := conDataRequest + conSData + "&loadtime=" + strconv.Itoa(poll.LoadTime)
	loadParams := "&dataversion=" + strconv.Itoa(poll.DataVersion) + "&timeout=60" + "&minimumdelay=60"
	// GET Request
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url+dataParams+loadParams, nil)
	if err != nil {
		return err
	}
	// Set Required Headers
	req.Header.Set("MMSSession", con.SessionToken)
	log.Println("Started Polling")

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	// Testing Polling data
	// log.Println(string(bodyBytes))
	if string(bodyBytes) == "" {
		return err
	}
	// log.Println("Polling Success")

	// Decode SData and add to struct
	sData := SData{}

	err = json.Unmarshal(bodyBytes, &sData)
	if err != nil {
		log.Println(string(bodyBytes))
		log.Println("Marshal bad")
		log.Println(err)

		return err
	}

	poll.Controller.m.Lock()
	defer poll.Controller.m.Unlock()

	if sData.Full == 1 { //nolint:nestif
		// Update all if data is full
		con.SData = sData
	} else {
		// Send con.Updated chan a message when a switch has been updated
		updated := false
		// Update device changes
		for _, d := range sData.Devices {
			// Update switch changes
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
			// Update lock changes
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

	// Update polling params read http://wiki.micasaverde.com/index.php/UI_Simple#lu_sdata:_The_polling_loop
	poll.DataVersion = sData.Dataversion
	poll.CurrentMinimumDelay = 2000
	poll.LoadTime = sData.Loadtime
	con.SData = sData

	return nil
}

// Close Removes controller from Vera and kills polling.
func (con *Controller) Close() error {
	con.m.Lock()
	con.Kill <- true
	// delete controller from vera identity
	err := con.Vera.removeDevice(con.DeviceID)
	if err != nil {
		return err
	}

	con.m.Unlock()

	return nil
}

// SwitchPowerStatus change swtich status.
func (con *Controller) SwitchPowerStatus(id int, status int) error {
	con.m.Lock()
	url := https + con.ServerRelay + conRelayPath + con.DeviceID
	dataParams := conDataRequest + conDevice + strconv.Itoa(id) + "&serviceId=" + conSwitch
	verParams := "&action=SetTarget&newTargetValue=" + strconv.Itoa(status)

	err := con.callURL(url + dataParams + verParams)
	if err != nil {
		return err
	}
	defer con.m.Unlock()

	return nil
}

// DoorLockStatus change lock status.
func (con *Controller) DoorLockStatus(id int, status int) error {
	con.m.Lock()
	url := https + con.ServerRelay + conRelayPath + con.DeviceID + conDataRequest
	serviceParams := conDevice + strconv.Itoa(id) + "&serviceId=" + conDoorLock
	actionParams := "&action=SetTarget&newTargetValue=" + strconv.Itoa(status)

	err := con.callURL(url + serviceParams + actionParams)
	if err != nil {
		return err
	}
	defer con.m.Unlock()

	return nil
}

// CustomRequest custom GET request controller using custom params
// params can be found http://wiki.micasaverde.com/index.php/Luup_Requests
// params after /port_3480/data_request? + "params"
// This function will not return GET data.
func (con *Controller) CustomRequest(params string) error {
	con.m.Lock()
	url := https + con.ServerRelay + conRelayPath + con.DeviceID
	params = conDataRequest + "?" + params

	err := con.callURL(url + params)
	if err != nil {
		return err
	}
	defer con.m.Unlock()

	return nil
}

// CustomRequestReturn custom GET request controller using custom params
// params can be found http://wiki.micasaverde.com/index.php/Luup_Requests
// params after /port_3480/data_request? + "params"
// This function will return back struct.
func (con *Controller) CustomRequestReturn(params string, data interface{}) (interface{}, error) {
	con.m.Lock()
	url := https + con.ServerRelay + conRelayPath + con.DeviceID
	params = conDataRequest + "?" + params

	data, err := con.callURLReturn(url+params, data)
	if err != nil {
		return data, err
	}
	defer con.m.Unlock()

	return data, nil
}

// callURL using GET.
func (con *Controller) callURL(url string) error {
	// GET Request
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	// Set Required Headers
	req.Header.Set("MMSSession", con.SessionToken)

	r, err := client().Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return nil
}

// callURLReturn using GET.
func (con *Controller) callURLReturn(url string, target interface{}) (interface{}, error) {
	// GET Request
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return target, err
	}

	// Set Required Headers
	req.Header.Set("MMSSession", con.SessionToken)

	r, err := client().Do(req)
	if err != nil {
		return target, err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target), nil
}
