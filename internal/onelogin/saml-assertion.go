package onelogin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const MFA_REQUIRED_STRING = "MFA is required for this user"

type SAMLAssertionBody struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
	AppID           string `json:"app_id"`
	SubDomain       string `json:"subdomain"`
}
type SAMLAssertionResponse struct {
	StateToken  string   `json:"state_token"`
	Message     string   `json:"message"`
	Devices     []Device `json:"devices"`
	CallbackURL string   `json:"callback_url"`
	User        User     `json:"user"`
}
type Device struct {
	DeviceID   int    `json:"device_id"`
	DeviceType string `json:"device_type"`
}
type User struct {
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	ID        int    `json:"id"`
}

const SAMLAssertionURl = OneLoginAPIURL + "api/2/saml_assertion"

func SAMLAssertion(token, login, password, appID, oneloginDomain string) (*SAMLAssertionResponse, error) {

	newBody := SAMLAssertionBody{
		UsernameOrEmail: login,
		Password:        password,
		AppID:           appID,
		SubDomain:       oneloginDomain,
	}

	jsonBody, err := json.MarshalIndent(newBody, "", "")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", SAMLAssertionURl, bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer:"+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var responseObject SAMLAssertionResponse
	err = json.Unmarshal(respBody, &responseObject)
	if err != nil {
		return nil, err
	}
	if responseObject.Message != "Success" && responseObject.Message != MFA_REQUIRED_STRING {
		return nil, fmt.Errorf(responseObject.Message)
	}

	return &responseObject, nil
}
