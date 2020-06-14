package main

// Vera struct hold info about one user
type Vera struct {
	Username     string
	Password     string
	Identity     IdentityJSON // NOTE Identity expires in 24 hrs
	AccountID    string
	SessionToken string
	Devices      Devices
}

// VeraController struct hold info about one controller
type VeraController struct {
	DeviceID     string
	ServerRelay  string
	SessionToken string
	SData        SData
	Switches     []Switch
}

//Switch devices with ON/OFF from vera controller
type Switch struct {
	ID     int    `json:"ID"`
	Name   string `json:"Name"`
	Status string `json:"Status"`
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

//SData struct to store data
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
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"sections"`
	Rooms []struct {
		Name    string `json:"name"`
		ID      int    `json:"id"`
		Section int    `json:"section"`
	} `json:"rooms"`
	Scenes []struct {
		Active int    `json:"active"`
		Name   string `json:"name"`
		ID     int    `json:"id"`
		Room   int    `json:"room"`
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

//SDataDevice struct for devices in SData
type SDataDevice struct {
	Name         string `json:"name"`
	Altid        string `json:"altid"`
	ID           int    `json:"id"`
	Category     int    `json:"category"`
	Subcategory  int    `json:"subcategory"`
	Room         int    `json:"room"`
	Parent       int    `json:"parent"`
	Configured   string `json:"configured"`
	State        int    `json:"state,omitempty"`
	Comment      string `json:"comment,omitempty"`
	Status       string `json:"status,omitempty"`
	Kwh          string `json:"kwh,omitempty"`
	Watts        string `json:"watts,omitempty"`
	Locked       string `json:"locked,omitempty"`
	Pincodes     string `json:"pincodes,omitempty"`
	CommFailure  string `json:"commFailure,omitempty"`
	Batterylevel string `json:"batterylevel,omitempty"`
}

//SDataCategory struct for category in SData
type SDataCategory struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}
