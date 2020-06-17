package govera

import "log"

//New Create new Vera object
func New(username string, password string) Vera {
	//Initialise Object
	vera := Vera{
		Username:    username,
		Password:    password,
		Controllers: &[]VeraController{},
	}

	//Gets Identity using username and password
	err := vera.GetIdentityToken()
	if err != nil {
		log.Panic(err)
	}
	//Gets SessionToken using Identity
	err = vera.GetSessionToken()
	if err != nil {
		log.Panic(err)
	}
	//Gets all devices linked to account using SessioToken
	err = vera.GetAllDevices()
	if err != nil {
		log.Panic(err)
	}

	return vera
}

//Renew Used to renew identity
func (vera *Vera) Renew() {

}
