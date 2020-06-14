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
