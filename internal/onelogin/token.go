package onelogin

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const OneLoginAPIURL = "https://api.us.onelogin.com/"

type AccessTokenResponse struct {
	Status struct {
		Error   bool   `json:"error"`
		Code    int    `json:"code"`
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"status"`
	Data []struct {
		AccessToken  string    `json:"access_token"`
		CreatedAt    time.Time `json:"created_at"`
		ExpiresIn    int       `json:"expires_in"`
		RefreshToken string    `json:"refresh_token"`
		TokenType    string    `json:"token_type"`
		AccountId    int       `json:"account_id"`
	} `json:"data"`
}

func GetAccessToken(id, secret string) (string, error) {

	body := "{\"grant_type\":\"client_credentials\"}"

	req, err := http.NewRequest("POST", OneLoginAPIURL+"auth/oauth2/token", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "client_id:"+id+", client_secret:"+secret)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	var responseObject AccessTokenResponse
	err = json.Unmarshal(respBody, &responseObject)
	if err != nil {
		return "", err
	}
	if !responseObject.Status.Error {
		return responseObject.Data[0].AccessToken, nil
	}

	return "", nil
}
