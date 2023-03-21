package onelogin

type Iface interface {
	SAMLAssertion(client HttpClient, token, login, password, appID, oneloginDomain string) (*SAMLAssertionResponse, error)
	VerifyFactor(client HttpClient, token string, deviceID int, appID string, stateToken string, mfaCode string) (*VerifyFactorResponse, error)
	GetAccessToken(client HttpClient, id, secret string) (string, error)
}

type Client struct {
}
