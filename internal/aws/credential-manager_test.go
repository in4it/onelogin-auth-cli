package aws

import (
	"testing"
)

type MockCredentialFileGetter struct {
	CredentialOutput string
	CredentialExists bool
}
type MockCredentialFileWriter struct {
	Out []byte
}

func (c MockCredentialFileGetter) SetHomeDir(homedir string) {

}
func (c MockCredentialFileGetter) Get() (string, error) {
	return c.CredentialOutput, nil
}
func (c MockCredentialFileGetter) Exists() (bool, error) {
	return c.CredentialExists, nil
}

func (c MockCredentialFileWriter) SetHomeDir(homedir string) {

}
func (c *MockCredentialFileWriter) Write(b []byte) error {
	c.Out = b
	return nil
}

func TestSetCredentialsEmpty(t *testing.T) {
	mockGetter := MockCredentialFileGetter{
		CredentialOutput: "",
		CredentialExists: true,
	}
	mockWriter := MockCredentialFileWriter{}
	err := Client{}.SetCredentials(
		mockGetter,
		&mockWriter,
		"accessKey",
		"secretKey",
		"sessionToken",
		"us-eas-1",
		"myProfile",
	)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if string(mockWriter.Out) != formatCredential("accessKey", "secretKey", "sessionToken", "us-eas-1", "myProfile") {
		t.Errorf("Wrong output. Got: %s", mockWriter.Out)
	}
}

func TestSetCredentialsNonExisting(t *testing.T) {
	mockGetter := MockCredentialFileGetter{
		CredentialOutput: "",
		CredentialExists: false,
	}
	mockWriter := MockCredentialFileWriter{}
	err := Client{}.SetCredentials(
		mockGetter,
		&mockWriter,
		"accessKey",
		"secretKey",
		"sessionToken",
		"us-eas-1",
		"myProfile",
	)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if string(mockWriter.Out) != formatCredential("accessKey", "secretKey", "sessionToken", "us-eas-1", "myProfile") {
		t.Errorf("Wrong output. Got: %s", mockWriter.Out)
	}
}
func TestSetCredentialsExistingFile(t *testing.T) {
	mockGetter := MockCredentialFileGetter{
		CredentialOutput: formatCredential("accessKey1", "secretKey1", "sessionToken1", "us-east-2", "profile123"),
		CredentialExists: true,
	}
	mockWriter := MockCredentialFileWriter{}
	err := Client{}.SetCredentials(
		mockGetter,
		&mockWriter,
		"accessKey",
		"secretKey",
		"sessionToken",
		"us-east-1",
		"myProfile",
	)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	expected := formatCredential("accessKey1", "secretKey1", "sessionToken1", "us-east-2", "profile123") + "\n\n" + formatCredential("accessKey", "secretKey", "sessionToken", "us-east-1", "myProfile")
	if string(mockWriter.Out) != expected {
		t.Errorf("Wrong output. Got: %s", mockWriter.Out)
	}
}

func TestSetCredentialsExistingFileReplace1(t *testing.T) {
	mockGetter := MockCredentialFileGetter{
		CredentialOutput: formatCredential("accessKey1", "secretKey1", "sessionToken1", "us-east-2", "profile123") + "\n\n" + formatCredential("accessKey2", "secretKey2", "sessionToken2", "us-east-2", "myProfile"),
		CredentialExists: true,
	}
	mockWriter := MockCredentialFileWriter{}
	err := Client{}.SetCredentials(
		mockGetter,
		&mockWriter,
		"accessKey",
		"secretKey",
		"sessionToken",
		"us-east-1",
		"myProfile",
	)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	expected := formatCredential("accessKey1", "secretKey1", "sessionToken1", "us-east-2", "profile123") + "\n\n" + formatCredential("accessKey", "secretKey", "sessionToken", "us-east-1", "myProfile")
	if string(mockWriter.Out) != expected {
		t.Errorf("Wrong output. Got: %s", mockWriter.Out)
	}
}

func TestSetCredentialsExistingFileReplace2(t *testing.T) {
	mockGetter := MockCredentialFileGetter{
		CredentialOutput: formatCredential("accessKey2", "secretKey2", "sessionToken2", "us-east-2", "myProfile") + "\n\n" + formatCredential("accessKey1", "secretKey1", "sessionToken1", "us-east-2", "profile123"),
		CredentialExists: true,
	}
	mockWriter := MockCredentialFileWriter{}
	err := Client{}.SetCredentials(
		mockGetter,
		&mockWriter,
		"accessKey",
		"secretKey",
		"sessionToken",
		"us-east-1",
		"myProfile",
	)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	expected := formatCredential("accessKey", "secretKey", "sessionToken", "us-east-1", "myProfile") + "\n\n" + formatCredential("accessKey1", "secretKey1", "sessionToken1", "us-east-2", "profile123")
	if string(mockWriter.Out) != expected {
		t.Errorf("Wrong output. Got: %s\n\nExpected: %s", mockWriter.Out, expected)
	}
}

func TestSetCredentialsExistingFileReplace3(t *testing.T) {
	mockGetter := MockCredentialFileGetter{
		CredentialOutput: formatCredential("accessKey1", "secretKey1", "sessionToken1", "us-east-2", "profile123") + "\n\n" + formatCredential("accessKey2", "secretKey2", "sessionToken2", "us-east-2", "myProfile") + "\n\n" + formatCredential("accessKey1", "secretKey1", "sessionToken1", "us-east-2", "profile1234"),
		CredentialExists: true,
	}
	mockWriter := MockCredentialFileWriter{}
	err := Client{}.SetCredentials(
		mockGetter,
		&mockWriter,
		"accessKey",
		"secretKey",
		"sessionToken",
		"us-east-1",
		"myProfile",
	)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	expected := formatCredential("accessKey1", "secretKey1", "sessionToken1", "us-east-2", "profile123") + "\n\n" + formatCredential("accessKey", "secretKey", "sessionToken", "us-east-1", "myProfile") + "\n\n" + formatCredential("accessKey1", "secretKey1", "sessionToken1", "us-east-2", "profile1234")
	if string(mockWriter.Out) != expected {
		t.Errorf("Wrong output. Got: %s\n\nExpected: %s", mockWriter.Out, expected)
	}
}
