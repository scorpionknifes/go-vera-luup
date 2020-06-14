package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
	log.Println(url)
	log.Println(con.SessionToken)
	con.SData = sData

	switchID := -1
	switches := []Switch{}
	for _, category := range sData.Categories {
		if category.Name == "On/Off Switch" {
			switchID = category.ID
		}
	}
	if switchID == -1 {
		return nil
	}
	for _, device := range sData.Devices {
		if device.Category == switchID {
			switchDevice := Switch{Name: device.Name, ID: device.ID, Status: device.Status}
			switches = append(switches, switchDevice)
		}
	}
	con.Switches = switches
	return nil
}

//Polling loop to CheckStatus using http://wiki.micasaverde.com/index.php/UI_Simple#lu_sdata:_The_polling_loop
func (con *VeraController) Polling() {
	go func() {
		for {
			err := con.CheckStatus()
			if err != nil {
				log.Panic(err)
			}
		}
	}()
}

//CheckStatus used to recheck switch status
func (con *VeraController) CheckStatus() error {
	//Get Url
	url := https + con.ServerRelay + conRelayPath + con.DeviceID + conDataRequest //+ "id=lu_sdata&loadtime=" + loadTime + "&dataversion=" + dataVersion + "&timeout=60" + "&minimumdelay=" + currentMinimumDelay
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
	con.SData = sData
	return nil
}
