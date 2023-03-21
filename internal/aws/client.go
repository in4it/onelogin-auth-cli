package aws

import "github.com/aws/aws-sdk-go/service/sts"

type Client struct {
}

type ClientIface interface {
	SetCredentials(credFileGetter CredentialFileGetterIface, credFileWriter CredentialFileWriterIface, accessKey, secretAccessKey, sessionToken, region, profileName string) error
	AssumeRoleWithSAML(accountID string, role string, assertionPayload string, durationSeconds int64) (*sts.AssumeRoleWithSAMLOutput, error)
}
