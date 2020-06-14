package main

//New Create new Vera object
func New(username string, password string) Vera {
	//Initialise Object
	vera := Vera{Username: username, Password: password}

	//Gets Identity using username and password
	vera.GetIdentityToken()
	//Gets SessionToken using Identity
	vera.GetSessionToken()
	//Gets all devices linked to account using SessioToken
	vera.GetAllDevices()

	return vera
}
