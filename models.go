package vera

import (
	"encoding/json"
	"sync"
)

// Vera class struct hold info about one user
// Vera object has reference to all Controller childs it creates
type Vera struct {
	Username     string
	Password     string
	Identity     IdentityJSON // NOTE Identity expires in 24 hrs
	AccountID    string
	SessionToken string
	Devices      Devices
	Controllers  *[]Controller // Add link so we can renew Identity
	m            *sync.Mutex
}

// Controller class struct hold info about one controller
// Controller object has reference to parent Vera object
type Controller struct {
	Vera         *Vera
	DeviceID     string
	ServerRelay  string
	SessionToken string
	SData        SData     // SData will not be polled unlike Switches
	Switches     *[]Switch // Data here would be update to date
	Locks        *[]Lock
	Kill         chan bool
	Updated      chan bool
	m            *sync.Mutex
}

// Polling struct to poll Controller
// Polling is a class that is only used for polling in goroutine
type Polling struct {
	LoadTime            int
	DataVersion         int
	CurrentMinimumDelay int
	Controller          *Controller
}

// Switch devices with ON/OFF from vera controller
type Switch struct {
	ID     int    `json:"ID"`
	Name   string `json:"Name"`
	Status string `json:"Status"`
}

// Lock devices with lock/unlock from vera controller
type Lock struct {
	ID     int    `json:"ID"`
	Name   string `json:"Name"`
	Locked string `json:"Locked"`
}

// IdentityJSON parse GetLoginToken
type IdentityJSON struct {
	Identity          string `json:"Identity"`
	IdentitySignature string `json:"IdentitySignature"`
	ServerEvent       string `json:"Server_Event"`
	ServerEventAlt    string `json:"Server_Event_Alt"`
	ServerAccount     string `json:"Server_Account"`
	ServerAccountAlt  string `json:"Server_Account_Alt"`
}

// Devices list all devices linked to an account
type Devices struct {
	Devices []Device `json:"Devices"`
}

// Device represents on device (smart controler)
type Device struct {
	PKDevice        string `json:"PK_Device"`
	PKDeviceType    string `json:"PK_DeviceType"`
	PKDeviceSubType string `json:"PK_DeviceSubType"`
	MacAddress      string `json:"MacAddress"`
	ServerDevice    string `json:"Server_Device"`
	ServerDeviceAlt string `json:"Server_Device_Alt"`
	PKInstallation  string `json:"PK_Installation"`
	Name            string `json:"Name"`
	Using2G         string `json:"Using_2G"`
	DeviceAssigned  string `json:"DeviceAssigned"`
	Blocked         int    `json:"Blocked"`
}

// DeviceInfo for one device (smart controler)
type DeviceInfo struct {
	PKDevice             string `json:"PK_Device"`
	ServerRelay          string `json:"Server_Relay"`
	MacAddress           string `json:"MacAddress"`
	Using2G              int    `json:"Using_2G"`
	ExternalIP           string `json:"ExternalIP"`
	AccessiblePort       string `json:"AccessiblePort"`
	InternalIP           string `json:"InternalIP"`
	AliveDate            string `json:"AliveDate"`
	FirmwareVersion      string `json:"FirmwareVersion"`
	PriorFirmwareVersion string `json:"PriorFirmwareVersion"`
	UpgradeDate          string `json:"UpgradeDate"`
	Uptime               string `json:"Uptime"`
	ServerDevice         string `json:"Server_Device"`
	ServerEvent          string `json:"Server_Event"`
	ServerSupport        string `json:"Server_Support"`
	ServerStorage        string `json:"Server_Storage"`
	WifiSsid             string `json:"WifiSsid"`
	Timezone             string `json:"Timezone"`
	LocalPort            string `json:"LocalPort"`
	ZWaveLocale          string `json:"ZWaveLocale"`
	ZWaveVersion         string `json:"ZWaveVersion"`
	FKBranding           string `json:"FK_Branding"`
	Platform             string `json:"Platform"`
	UILanguage           string `json:"UILanguage"`
	UISkin               string `json:"UISkin"`
	HasWifi              string `json:"HasWifi"`
	HasAlarmPanel        string `json:"HasAlarmPanel"`
	UI                   string `json:"UI"`
	EngineStatus         string `json:"EngineStatus"`
	DistributionBuild    string `json:"DistributionBuild"`
	AccessPermissions    string `json:"AccessPermissions"`
	LinuxFirmware        int    `json:"LinuxFirmware"`
}

// SData struct to store SData from Luup Request on a device
// SData struct is designed to be public and users can store data in alternative locations
// json.Number is used cause Luup Request is not consistent when polling
// e.g int is received when full=1 and string is received when full=0
type SData struct {
	Full         int    `json:"full"`
	Version      string `json:"version"`
	Model        string `json:"model"`
	ZwaveHeal    int    `json:"zwave_heal"`
	Temperature  string `json:"temperature"`
	Skin         string `json:"skin"`
	SerialNumber string `json:"serial_number"`
	Fwd1         string `json:"fwd1"`
	Fwd2         string `json:"fwd2"`
	Mode         int    `json:"mode"`
	Sections     []struct {
		Name string      `json:"name"`
		ID   json.Number `json:"id"`
	} `json:"sections"`
	Rooms []struct {
		Name    string      `json:"name"`
		ID      json.Number `json:"id"`
		Section json.Number `json:"section"`
	} `json:"rooms"`
	Scenes []struct {
		Name    string      `json:"name"`
		ID      json.Number `json:"id"`
		Room    json.Number `json:"room"`
		Active  json.Number `json:"active"`
		State   json.Number `json:"state"`
		Comment string      `json:"comment"`
	} `json:"scenes"`
	Devices     []SDataDevice   `json:"devices"`
	Categories  []SDataCategory `json:"categories"`
	Ir          int             `json:"ir"`
	Irtx        string          `json:"irtx"`
	Loadtime    int             `json:"loadtime"`
	Dataversion int             `json:"dataversion"`
	State       int             `json:"state"`
	Comment     string          `json:"comment"`
}

// SDataDevice struct for devices in SData
type SDataDevice struct {
	Name         string      `json:"name"`
	Altid        string      `json:"altid"`
	ID           json.Number `json:"id"` // SData when polling returns string instead of int
	Category     json.Number `json:"category"`
	Subcategory  json.Number `json:"subcategory"`
	Room         json.Number `json:"room"`
	Parent       json.Number `json:"parent"`
	Configured   string      `json:"configured"`
	State        json.Number `json:"state"`
	Comment      string      `json:"comment"`
	Kwh          string      `json:"kwh,omitempty"`
	Status       string      `json:"status"`
	Watts        string      `json:"watts,omitempty"`
	Pincodes     string      `json:"pincodes,omitempty"`
	CommFailure  string      `json:"commFailure,omitempty"`
	Batterylevel string      `json:"batterylevel,omitempty"`
	Locked       string      `json:"locked,omitempty"`
}

// SDataCategory struct for category in SData
type SDataCategory struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}
