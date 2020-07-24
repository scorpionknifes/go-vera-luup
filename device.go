package vera

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

// GetAllDevices linked to account
// Get new list of devices linked to vera account
// This should be used after adding new devices tp a live deployment
func (vera *Vera) GetAllDevices() error {
	vera.m.Lock()
	defer vera.m.Unlock()
	url := https + vera.Identity.ServerAccount + accountPath + vera.AccountID + devicesPath
	err := vera.getAllDevicesURL(url)
	if err == nil {
		return nil
	}
	// if error occurred try using ServerAccountAlt
	url = https + vera.Identity.ServerAccountAlt + accountPath + vera.AccountID + devicesPath
	return vera.getAllDevicesURL(url)
}

func (vera *Vera) getAllDevicesURL(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	// Set Required Headers
	req.Header.Set("MMSSession", vera.SessionToken)
	r, err := client.Do(req)
	if err != nil {
		return err
	}
	// Decode devices and add to Vera struct
	devices := Devices{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&devices)
	if err != nil {
		return err
	}
	vera.Devices = devices
	return nil
}

// GetDeviceRelay construct a Controller Object for the deviceID given
// Each new Vera controller would require a new Object and SessionToken to interact
// Each SessionToken also expires in 24 hrs the Object would auto renew with Vera Identity
func (vera *Vera) GetDeviceRelay(deviceID string) (Controller, error) {
	vera.m.Lock()
	deviceInfo, err := vera.getDeviceInfo(deviceID)
	if err != nil {
		return Controller{}, err
	}
	controller := Controller{
		Vera:        vera,
		DeviceID:    deviceID,
		ServerRelay: deviceInfo.ServerRelay,
		Kill:        make(chan bool),
		Updated:     make(chan bool),
		m:           &sync.Mutex{},
	}
	err = controller.GetSessionToken()
	if err != nil {
		return Controller{}, err
	}
	err = controller.GetSData()
	if err != nil {
		return Controller{}, err
	}
	// Enable Polling using go routine
	controller.Polling()

	// Add controller to controller array in Vera
	*vera.Controllers = append(*vera.Controllers, controller)
	vera.m.Unlock()
	return controller, err
}

// getDeviceInfo get device info of deviceID
func (vera *Vera) getDeviceInfo(deviceID string) (DeviceInfo, error) {
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
	deviceInfo, err := vera.getDeviceInfoURL(url)

	// Try using ServerDeviceAlt if ServerDevice doesn't work
	if err != nil {
		url = https + device.ServerDeviceAlt + devicePath + deviceID
		deviceInfo, err = vera.getDeviceInfoURL(url)
	}
	return deviceInfo, err
}

func (vera *Vera) getDeviceInfoURL(url string) (DeviceInfo, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return DeviceInfo{}, err
	}
	// Set Required Headers
	req.Header.Set("MMSSession", vera.SessionToken)
	r, err := client.Do(req)
	if err != nil {
		return DeviceInfo{}, err
	}
	// Decode devices and add to Vera struct
	deviceInfo := DeviceInfo{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&deviceInfo)
	if err != nil {
		return DeviceInfo{}, err
	}
	return deviceInfo, nil
}

// removeDevice remove controller device from Vera account
func (vera *Vera) removeDevice(deviceID string) error {
	cons := *vera.Controllers
	for i, device := range cons {
		if device.DeviceID == deviceID {
			cons[i] = cons[len(cons)-1]
			*vera.Controllers = cons[:len(cons)-1]
			return nil
		}
	}
	return errors.New("DeviceID not found")
}
