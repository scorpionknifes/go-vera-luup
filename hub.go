package main

import (
	"log"
	"net/http"
)

const (
	account = "/account/account/account/"
	devices = "/devices"
)

//GetAllDevices linked to account
func (vera *Vera) GetAllDevices() error {
	url := https + vera.Identity.ServerAccount + account + vera.AccountID + devices
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	//Set Required Headers
	req.Header.Set("MMSSession", vera.SessionToken)
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	log.Println(vera.SessionToken)
	log.Println(url)
	return nil
}
