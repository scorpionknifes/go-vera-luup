package main

import (
	"encoding/json"
	"net/http"
)

const (
	accountPath = "/account/account/account/"
	devicesPath = "/devices"
)

//GetAllDevices linked to account
func (vera *Vera) GetAllDevices() error {
	url := https + vera.Identity.ServerAccount + accountPath + vera.AccountID + devicesPath
	err := vera.GetAllDevicesURL(url)
	if err == nil {
		return nil
	}
	//if error occured try using ServerAccountAlt
	url = https + vera.Identity.ServerAccountAlt + accountPath + vera.AccountID + devicesPath
	return vera.GetAllDevicesURL(url)
}

//GetAllDevicesURL linked to account using URL
func (vera *Vera) GetAllDevicesURL(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	//Set Required Headers
	req.Header.Set("MMSSession", vera.SessionToken)
	r, err := client.Do(req)
	if err != nil {
		return err
	}
	//Decode devices and add to Vera struct
	devices := Devices{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&devices)
	if err != nil {
		return err
	}
	vera.Devices = devices
	return nil
}
