package onelogin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type VerifyFactorBody struct {
	AppID      string `json:"app_id"`
	DeviceID   string `json:"device_id"`
	StateToken string `json:"state_token"`
	OTPToken   string `json:"otp_token"`
}

type VerifyFactorResponse struct {
	Data   string `json:"data"`
	Status struct {
		Message string `json:"message"`
		Error   bool   `json:"error"`
		Type    string `json:"type"`
		Code    int    `json:"code"`
	} `json:"status"`
}

const VerifyFactorURL = OneLoginAPIURL + "api/1/saml_assertion/verify_factor"

func VerifyFactor(token string, deviceID int, appID string, stateToken string, mfaCode string) (*VerifyFactorResponse, error) {
	newBody := VerifyFactorBody{
		AppID:      appID,
		DeviceID:   strconv.Itoa(deviceID),
		StateToken: stateToken,
		OTPToken:   mfaCode,
	}

	jsonBody, err := json.MarshalIndent(newBody, "", "")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", VerifyFactorURL, bytes.NewBuffer([]byte(jsonBody)))
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
	var responseObject VerifyFactorResponse
	err = json.Unmarshal(respBody, &responseObject)
	if err != nil {
		return nil, err
	}
	if responseObject.Status.Code != 200 {
		return nil, fmt.Errorf(responseObject.Status.Message)
	}

	return &responseObject, nil
}
