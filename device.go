package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	devicePath = "/device/device/device/"
)

//GetDeviceInfo get device info of deviceID
func (vera *Vera) GetDeviceInfo(deviceID string) error {
	var device Device
	for _, d := range vera.Devices.Devices {
		if deviceID == d.PKDevice {
			device = d
		}
	}
	if device == (Device{}) {
		return errors.New("deviceID '" + deviceID + "' not found")
	}
	url := https + device.ServerDevice + devicePath + deviceID
	log.Println(url)
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
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	log.Println(string(bodyBytes))
	return nil
}
