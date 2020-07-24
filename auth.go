package vera

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/buger/jsonparser"
)

// GetIdentityToken gets the identity token from vera using username and password
func (vera *Vera) GetIdentityToken() error {
	//Get Url
	url := https + remoteURL + loginPath + vera.Username
	//Get Url Params
	hashedpassword := hash(vera.Username + vera.Password + passwordSeed)
	params := "?SHA1Password=" + hashedpassword + "&PK_Oem=1"
	//GET Request
	identity := IdentityJSON{}
	err := getJSON(url+params, &identity)
	if err != nil {
		return err
	}
	vera.Identity = identity

	//Get PK_Account aka AccountID by decode IdentityToken
	raw, err := base64.StdEncoding.DecodeString(identity.Identity)
	if err != nil {
		return err
	}
	accountID, err := jsonparser.GetInt(raw, "PK_Account")
	if err != nil {
		return err
	}
	vera.AccountID = strconv.FormatInt(accountID, 10)
	return err
}

// GetSessionToken gets the session token using identity token
func (vera *Vera) GetSessionToken() error {
	url := https + vera.Identity.ServerAccount + sessionPath
	err := vera.GetSessionTokenURL(url)
	if err == nil {
		return nil
	}
	//if error occurred try using ServerAccountAlt
	url = https + vera.Identity.ServerAccountAlt + sessionPath
	return vera.GetSessionTokenURL(url)
}

// GetSessionTokenURL gets the session token using identity token and URL
func (vera *Vera) GetSessionTokenURL(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	//Set Required Headers
	req.Header.Set("MMSAuth", vera.Identity.Identity)
	req.Header.Set("MMSAuthSig", vera.Identity.IdentitySignature)
	r, err := client.Do(req)
	if err != nil {
		return err
	}
	//Convert response into string as SessionToken
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	vera.SessionToken = string(bodyBytes)
	return nil
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
