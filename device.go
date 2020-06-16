package govera

import (
	"encoding/json"
	"errors"
	"net/http"
)

const (
	devicePath = "/device/device/device/"
)

//GetDeviceRelay get device relay
func (vera *Vera) GetDeviceRelay(deviceID string) (VeraController, error) {
	deviceInfo, err := vera.GetDeviceInfo(deviceID)
	if err != nil {
		return VeraController{}, err
	}
	controller := VeraController{
		DeviceID:    deviceID,
		ServerRelay: deviceInfo.ServerRelay,
		Kill:        make(chan bool),
		Updated:     make(chan bool),
	}
	err = controller.GetSessionToken(vera.Identity)
	if err != nil {
		return VeraController{}, err
	}
	err = controller.GetSData()
	if err != nil {
		return VeraController{}, err
	}
	//Enable Polling using go routine
	controller.Polling()

	return controller, err
}

//GetDeviceInfo get device info of deviceID
func (vera *Vera) GetDeviceInfo(deviceID string) (DeviceInfo, error) {
	var device Device
	for _, d := range vera.Devices.Devices {
		if deviceID == d.PKDevice {
			device = d
		}
	}
	if device == (Device{}) {
		return DeviceInfo{}, errors.New("deviceID '" + deviceID + "' not found")
	}
	url := https + device.ServerDevice + devicePath + deviceID
	deviceInfo, err := vera.GetDeviceInfoURL(url)

	//Try using ServerDeviceAlt if ServerDevice doesn't work
	if err != nil {
		url = https + device.ServerDeviceAlt + devicePath + deviceID
		deviceInfo, err = vera.GetDeviceInfoURL(url)
	}
	return deviceInfo, err
}

//GetDeviceInfoURL get device info using url
func (vera *Vera) GetDeviceInfoURL(url string) (DeviceInfo, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return DeviceInfo{}, err
	}
	//Set Required Headers
	req.Header.Set("MMSSession", vera.SessionToken)
	r, err := client.Do(req)
	if err != nil {
		return DeviceInfo{}, err
	}
	//Decode devices and add to Vera struct
	deviceInfo := DeviceInfo{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&deviceInfo)
	if err != nil {
		return DeviceInfo{}, err
	}
	return deviceInfo, nil
}
