package main

// Vera struct hold info about one controller
type Vera struct {
	SerialNumber string
	Username     string
	Password     string
	Identity     IdentityJSON // NOTE Identity expires in 24 hrs
	AccountID    string
	SessionToken string
	Devices      Devices
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

// Device represents on device
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
