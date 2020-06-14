package main

// Vera struct hold info about one controller
type Vera struct {
	SerialNumber string
	Username     string
	Password     string
	Identity     IdentityJSON
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
