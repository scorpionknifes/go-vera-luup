package vera

import ( //nolint:gci
	"context"
	"crypto/sha1" //nolint:gosec
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/buger/jsonparser"
)

// GetIdentityToken gets the identity token from vera using username and password.
func (vera *Vera) GetIdentityToken() error {
	// Get Url
	url := https + remoteURL + loginPath + vera.Username

	// Get Url Params
	hashedpassword, err := hash(vera.Username + vera.Password + passwordSeed)
	if err != nil {
		return err
	}

	params := "?SHA1Password=" + hashedpassword + "&PK_Oem=1"
	// GET Request
	vera.Identity = IdentityJSON{}

	if err = getJSON(url+params, &vera.Identity); err != nil {
		return err
	}

	// Get PK_Account aka AccountID by decode IdentityToken
	raw, err := base64.StdEncoding.DecodeString(vera.Identity.Identity)
	if err != nil {
		return err
	}

	accountID, err := jsonparser.GetInt(raw, "PK_Account")
	if err != nil {
		return err
	}

	vera.AccountID = strconv.FormatInt(accountID, 10)

	return nil
}

// GetSessionToken gets the session token using identity token.
func (vera *Vera) GetSessionToken() error {
	url := https + vera.Identity.ServerAccount + sessionPath

	err := vera.GetSessionTokenURL(url)
	if err == nil {
		return nil
	}

	// if error occurred try using ServerAccountAlt
	url = https + vera.Identity.ServerAccountAlt + sessionPath

	return vera.GetSessionTokenURL(url)
}

// GetSessionTokenURL gets the session token using identity token and URL.
func (vera *Vera) GetSessionTokenURL(url string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	// Set Required Headers
	req.Header.Set("MMSAuth", vera.Identity.Identity)
	req.Header.Set("MMSAuthSig", vera.Identity.IdentitySignature)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// Convert response into string as SessionToken
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	vera.SessionToken = string(bodyBytes)

	return nil
}

// hash using sha1.
func hash(s string) (string, error) {
	h := sha1.New() //nolint:gosec
	_, err := h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil)), err
}

// getJSON get json by using url.
func getJSON(url string, target interface{}) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	return json.NewDecoder(req.Body).Decode(target)
}
