package main

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

	switchID := ""
	switches := []Switch{}
	for _, category := range sData.Categories {
		if category.Name == "On/Off Switch" {
			switchID = string(category.ID)
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
	}
	con.Switches = switches
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
					log.Panic(err)
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
		log.Fatal(err)
	}
	//Testing Polling data
	log.Println(string(bodyBytes))
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
		//Update switch changes
		for _, d := range sData.Devices {
			for i, e := range con.Switches {
				id, _ := strconv.Atoi(string(d.ID))
				if id == e.ID {
					e.Status = d.Status
					con.Switches[i] = e
					break
				}
			}
		}
	}

	//Update polling params read http://wiki.micasaverde.com/index.php/UI_Simple#lu_sdata:_The_polling_loop
	poll.DataVersion = sData.Dataversion
	poll.CurrentMinimumDelay = 2000
	poll.LoadTime = sData.Loadtime
	con.SData = sData
	return nil
}
