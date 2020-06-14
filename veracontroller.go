package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	conSessionPath = "/info/session/token"
	conRelayPath   = "/relay/relay/relay/device/"
	conPortPath    = "/port_3480/data_request?id=sdata"
)

//GetSessionToken get relay session by using identity
func (con *VeraController) GetSessionToken(identity IdentityJSON) error {
	//Get Url
	url := https + con.ServerRelay + conSessionPath
	//GET Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	//Set Required Headers
	req.Header.Set("MMSAuth", identity.Identity)
	req.Header.Set("MMSAuthSig", identity.IdentitySignature)
	r, err := client.Do(req)
	if err != nil {
		return err
	}
	//Convert response into string as SessionToken
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	con.SessionToken = string(bodyBytes)
	return nil
}

//GetSData Get SData from Hub through Relay Server aka all info
func (con *VeraController) GetSData() error {
	//Get Url
	url := https + con.ServerRelay + conRelayPath + con.DeviceID + conPortPath
	//GET Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	//Set Required Headers
	req.Header.Set("MMSSession", con.SessionToken)
	r, err := client.Do(req)
	if err != nil {
		return err
	}
	//Decode SData and add to struct
	sData := SData{}
	err = json.NewDecoder(r.Body).Decode(&sData)
	if err != nil {
		return err
	}
	con.SData = sData
	return nil
}
