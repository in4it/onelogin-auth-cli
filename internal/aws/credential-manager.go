package aws

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type CredentialFileGetterIface interface {
	Get() (string, error)
	Exists() (bool, error)
	SetHomeDir(string)
}
type CredentialFileWriterIface interface {
	Write(b []byte) error
	SetHomeDir(string)
}

type CredentialFileGetter struct {
	Homedir string
}

func (c CredentialFileGetter) SetHomeDir(homedir string) {
	c.Homedir = homedir
}
func (c CredentialFileGetter) Get() (string, error) {
	credentials, err := ioutil.ReadFile(c.Homedir + "/.aws/credentials")
	if err != nil {
		return "", fmt.Errorf("Could not read file: %s", err)
	}

	return string(credentials), nil
}
func (c CredentialFileGetter) Exists() (bool, error) {
	credentialsFileExists := false

	if _, err := os.Stat(c.Homedir + "/.aws/credentials"); err == nil {
		credentialsFileExists = true
	} else if errors.Is(err, os.ErrNotExist) {
		credentialsFileExists = false
	} else {
		return credentialsFileExists, fmt.Errorf("Could not determine if credential file exists: %s", err)
	}
	return credentialsFileExists, nil
}

type CredentialFileWriter struct {
	Homedir string
}

func (c CredentialFileWriter) Write(b []byte) error {
	return os.WriteFile(c.Homedir+"/.aws/credentials", b, 0600)
}
func (c CredentialFileWriter) SetHomeDir(homedir string) {
	c.Homedir = homedir
}

func SetCredentials(credFileGetter CredentialFileGetterIface, credFileWriter CredentialFileWriterIface, accessKey, secretAccessKey, sessionToken, region, profileName string) error {
	// set user homedir
	homedir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user's home directory: %s", err)
	}
	credFileGetter.SetHomeDir(homedir)
	credFileWriter.SetHomeDir(homedir)

	// check whether file exists
	exists, err := credFileGetter.Exists()
	if err != nil {
		return fmt.Errorf("Couldn't determine whether credential file exists: %s", err)
	}

	// if exists, update / write file
	if exists {
		credentials, err := credFileGetter.Get()
		if err != nil {
			return fmt.Errorf("Couldn't get credential file: %s", err)
		}
		if profileExists(credentials, profileName) {
			err = updateCredential(credFileWriter, credentials, accessKey, secretAccessKey, sessionToken, region, profileName)
		} else {
			appendCredential(credFileWriter, credentials, accessKey, secretAccessKey, sessionToken, region, profileName)
		}
	} else { // write new file if credential file doesn't exist
		writeCredentialAsNewFile(credFileWriter, accessKey, secretAccessKey, sessionToken, region, profileName)
	}

	return nil
}

func profileExists(credentials string, profileName string) bool {
	scanner := bufio.NewScanner(strings.NewReader(credentials))
	for scanner.Scan() {
		if strings.Trim(scanner.Text(), " ") == "["+profileName+"]" {
			return true
		}
	}
	return false
}

func updateCredential(writer CredentialFileWriterIface, credentials, accessKey, secretAccessKey, sessionToken, region, profileName string) error {
	newCredentials := ""
	found := false
	scanner := bufio.NewScanner(strings.NewReader(credentials))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Trim(line, " ") == "["+profileName+"]" {
			found = true
			newCredentials += formatCredential(accessKey, secretAccessKey, sessionToken, region, profileName)
		}

		if !found {
			newCredentials += line + "\n"
		}

		if found && line == "" {
			found = false
			newCredentials += "\n\n"
		}
	}
	newCredentials = strings.TrimRight(newCredentials, "\n")
	return writer.Write([]byte(newCredentials))
}
func appendCredential(writer CredentialFileWriterIface, credentials, accessKey, secretAccessKey, sessionToken, region, profileName string) error {
	if credentials != "" {
		credentials += "\n\n"
	}
	credentials += formatCredential(accessKey, secretAccessKey, sessionToken, region, profileName)
	return writer.Write([]byte(credentials))
}

func writeCredentialAsNewFile(writer CredentialFileWriterIface, accessKey, secretAccessKey, sessionToken, region, profileName string) error {
	return writer.Write([]byte(formatCredential(accessKey, secretAccessKey, sessionToken, region, profileName)))
}

func formatCredential(accessKey, secretAccessKey, sessionToken, region, profileName string) string {
	return fmt.Sprintf("[%s]\naws_access_key_id = %s\naws_secret_access_key = %s\naws_session_token = %s", profileName, accessKey, secretAccessKey, sessionToken)
}
