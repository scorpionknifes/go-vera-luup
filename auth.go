package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	remoteURL    = "https://us-autha11.mios.com"
	loginURL     = "/autha/auth/username/"
	passwordSeed = "oZ7QE6LcLJp6fiWzdqZc"
)

var client = &http.Client{Timeout: 10 * time.Second}

// GetLoginToken gets the login token from vera
func (vera *Vera) GetLoginToken() error {
	//Get Url
	url := remoteURL + loginURL + vera.Username
	//Get Url Params
	hashedpassword := hash(vera.Username + vera.Password + passwordSeed)
	params := "?SHA1Password=" + hashedpassword + "&PK_Oem=1"
	//GET Request
	log.Println(url + params)
	identity := IdentityJSON{}
	err := getJSON(url+params, &identity)
	vera.Identity = identity
	return err
}

//hash using sha1
func hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//getJSON get json by using url
func getJSON(url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
