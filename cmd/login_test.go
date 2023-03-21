package cmd

import (
	intAWS "onelogin-auth-cli/internal/aws"
	"onelogin-auth-cli/internal/onelogin"
	"onelogin-auth-cli/utils"

	"testing"

	"github.com/aws/aws-sdk-go/service/sts"
)

type promptMock struct {
	GetRoleResult       *int
	GetAccountResult    *int
	GetDeviceIDResult   *int
	GetEmailResponse    string
	GetPasswordResponse string
}

func (p promptMock) getRole(promptRunner utils.PromptRunner, roles []string) (*int, error) {
	return p.GetRoleResult, nil
}
func (p promptMock) getAccount(promptRunner utils.PromptRunner, accounts []Account) (*int, error) {
	return p.GetAccountResult, nil
}
func (p promptMock) getDeviceID(promptRunner utils.PromptRunner, devices []onelogin.Device) (*int, error) {
	return p.GetDeviceIDResult, nil
}
func (p promptMock) getEmail() (string, error) {
	return p.GetEmailResponse, nil
}
func (p promptMock) getPassword() (string, error) {
	return p.GetPasswordResponse, nil
}

type runPromptMock struct {
	RunResultInt    int
	RunResultString string
}

func (p runPromptMock) Run() (int, string, error) {
	return p.RunResultInt, p.RunResultString, nil
}

type mockOneloginClient struct {
	SAMLAssertationResponse *onelogin.SAMLAssertionResponse
	VerifyFactorResponse    *onelogin.VerifyFactorResponse
	GetAccessTokenResponse  string
}

func (m mockOneloginClient) SAMLAssertion(client onelogin.HttpClient, token, login, password, appID, oneloginDomain string) (*onelogin.SAMLAssertionResponse, error) {
	return m.SAMLAssertationResponse, nil
}
func (m mockOneloginClient) VerifyFactor(client onelogin.HttpClient, token string, deviceID int, appID string, stateToken string, mfaCode string) (*onelogin.VerifyFactorResponse, error) {
	return m.VerifyFactorResponse, nil
}
func (m mockOneloginClient) GetAccessToken(client onelogin.HttpClient, id, secret string) (string, error) {
	return m.GetAccessTokenResponse, nil
}

type mockAwsClient struct {
	AssumeRoleWithSAMLResponse *sts.AssumeRoleWithSAMLOutput
	AccessKey                  string
	SecretAccessKey            string
	SessionToken               string
}

func (m *mockAwsClient) SetCredentials(credFileGetter intAWS.CredentialFileGetterIface, credFileWriter intAWS.CredentialFileWriterIface, accessKey, secretAccessKey, sessionToken, region, profileName string) error {
	m.AccessKey = accessKey
	m.SecretAccessKey = secretAccessKey
	m.SessionToken = sessionToken
	return nil
}
func (m *mockAwsClient) AssumeRoleWithSAML(accountID string, role string, assertionPayload string, durationSeconds int64) (*sts.AssumeRoleWithSAMLOutput, error) {
	return m.AssumeRoleWithSAMLResponse, nil
}

func TestGetRole(t *testing.T) {
	options := []string{"a", "b", "c"}
	answer := 2
	res, err := prompts{}.getRole(runPromptMock{RunResultString: options[answer]}, options)
	if err != nil {
		t.Fatalf("getRole error: %s", err)
	}
	if *res != answer {
		t.Fatalf("Got wrong answer: %d", *res)
	}
}

func TestGetRoleEmpty(t *testing.T) {
	options := []string{}
	_, err := prompts{}.getRole(runPromptMock{RunResultString: ""}, options)
	if err.Error() != "role not found" {
		t.Fatalf("Expected role not found error, got nil")
	}
}

func TestGetRoleOneItem(t *testing.T) {
	options := []string{"a"}
	answer := 0
	res, err := prompts{}.getRole(runPromptMock{RunResultString: options[answer]}, options)
	if err != nil {
		t.Fatalf("getRole error: %s", err)
	}
	if *res != answer {
		t.Fatalf("Got wrong answer: %d", *res)
	}
}
func TestGetRoleWrongItem(t *testing.T) {
	options := []string{"a", "b", "c"}
	_, err := prompts{}.getRole(runPromptMock{RunResultString: "d"}, options)
	if err.Error() != "role not found" {
		t.Fatalf("Expected role not found error, got nil")
	}
}

func TestGetAccount(t *testing.T) {
	accounts := []Account{
		{
			Name: "a",
		},
		{
			Name: "b",
		},
		{
			Name: "c",
		},
	}
	answer := 1
	res, err := prompts{}.getAccount(runPromptMock{RunResultString: accounts[answer].Name}, accounts)
	if err != nil {
		t.Fatalf("getAccount error: %s", err)
	}
	if *res != answer {
		t.Fatalf("Got wrong answer: %d", *res)
	}
}

func TestGetDeviceID(t *testing.T) {
	devices := []onelogin.Device{
		{
			DeviceID:   1,
			DeviceType: "type1",
		},
		{
			DeviceID:   2,
			DeviceType: "type2",
		},
		{
			DeviceID:   3,
			DeviceType: "type3",
		},
	}
	answer := 2
	res, err := prompts{}.getDeviceID(runPromptMock{RunResultString: devices[answer].DeviceType}, devices)
	if err != nil {
		t.Fatalf("getDeviceID error: %s", err)
	}
	if *res != devices[answer].DeviceID {
		t.Fatalf("Got wrong answer: %d", *res)
	}
}

func TestLogin(t *testing.T) {
	err := LoadConfig("../internal/testdata/")
	if err != nil {
		t.Fatalf("LoadConfig error: %s", err)
	}
	mockClient := mockOneloginClient{
		SAMLAssertationResponse: &onelogin.SAMLAssertionResponse{},
		VerifyFactorResponse:    &onelogin.VerifyFactorResponse{},
		GetAccessTokenResponse:  "accessToken",
	}
	promptResults := 0
	pm := promptMock{
		GetRoleResult:       &promptResults,
		GetAccountResult:    &promptResults,
		GetDeviceIDResult:   &promptResults,
		GetEmailResponse:    "user@domain.inv",
		GetPasswordResponse: "password",
	}
	ok1 := "ok1"
	ok2 := "ok2"
	ok3 := "ok3"
	a := mockAwsClient{
		AssumeRoleWithSAMLResponse: &sts.AssumeRoleWithSAMLOutput{
			Credentials: &sts.Credentials{
				AccessKeyId:     &ok1,
				SecretAccessKey: &ok2,
				SessionToken:    &ok3,
			},
		},
	}
	err = doLogin(mockClient, pm, &a, []string{})
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if a.AccessKey != ok1 {
		t.Fatalf("Wrong access key: %s", a.AccessKey)
	}
	if a.SecretAccessKey != ok2 {
		t.Fatalf("Wrong secret access key: %s", a.SecretAccessKey)
	}
	if a.SessionToken != ok3 {
		t.Fatalf("Wrong session token: %s", a.SessionToken)
	}
}
