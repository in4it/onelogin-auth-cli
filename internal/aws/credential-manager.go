package aws

import (
	"fmt"
	"log"
	"os/exec"
)

func SetCredentials(accessKey, secretAccessKey, sessionToken, region, profileName string) {
	err := setProfileParameter("aws_access_key_id", accessKey, profileName)
	if err != nil {
		log.Fatalln("Error while writing credentials to file. ", err)
	}
	err = setProfileParameter("aws_secret_access_key", secretAccessKey, profileName)
	if err != nil {
		log.Fatalln("Error while writing credentials to file. ", err)
	}
	err = setProfileParameter("aws_session_token", sessionToken, profileName)
	if err != nil {
		log.Fatalln("Error while writing credentials to file. ", err)
	}
	err = setProfileParameter("region", region, profileName)
	if err != nil {
		log.Fatalln("Error while writing credentials to file. ", err)
	}
}

func setProfileParameter(parameter, value, profileName string) error {
	command := []string{
		"aws",
		"configure",
		"set",
		parameter,
		value,
		"--profile",
		profileName,
	}

	cmd := exec.Command(command[0], command[1], command[2], command[3], command[4], command[5], command[6])
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
