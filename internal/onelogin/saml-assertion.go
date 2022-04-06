package onelogin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SAMLAssertionBody struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
	AppID           string `json:"app_id"`
	SubDomain       string `json:"subdomain"`
}

type SAMLAssertionResponse struct {
	Data []struct {
		Devices     []Device `json:"devices"`
		CallbackUrl string   `json:"callback_url"`
		User        struct {
			Username  string `json:"username"`
			Email     string `json:"email"`
			Lastname  string `json:"lastname"`
			Id        int    `json:"id"`
			Firstname string `json:"firstname"`
		} `json:"user"`
		StateToken string `json:"state_token"`
	} `json:"data"`
	Status struct {
		Message string `json:"message"`
		Error   bool   `json:"error"`
		Type    string `json:"type"`
		Code    int    `json:"code"`
	} `json:"status"`
}

type Device struct {
	DeviceId   int    `json:"device_id"`
	DeviceType string `json:"device_type"`
}

const SAMLAssertionURl = OneLoginAPIURL + "api/1/saml_assertion"

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
	if responseObject.Status.Code != 200 {
		return nil, fmt.Errorf(responseObject.Status.Message)
	}

	return &responseObject, nil
}
