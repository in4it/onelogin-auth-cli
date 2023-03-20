package cmd

import (
	"fmt"
	"log"
	intAWS "onelogin-auth-cli/internal/aws"
	"onelogin-auth-cli/internal/onelogin"
	"onelogin-auth-cli/utils"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login command",
	Run: func(cmd *cobra.Command, args []string) {
		var role, account *int
		var err error
		var assertionPayload string

		//Get Role and Accounts from parameters or from keyboard input
		if len(args) != 2 {
			role, err = getRole(config.Roles)
			if err != nil {
				log.Fatalln(err)
			}
			account, err = getAccount(config.Accounts)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			roleNum, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Role must be a number")
				os.Exit(1)
			}
			accountNum, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Account must be a number")
				os.Exit(1)
			}
			if roleNum > len(config.Roles)-1 {
				fmt.Println("Invalid Role")
				os.Exit(1)
			}
			if accountNum > len(config.Accounts)-1 {
				fmt.Println("Invalid Account")
				os.Exit(1)
			}
			role = &roleNum
			account = &accountNum
			fmt.Println("Role: ", config.Roles[*role])
			fmt.Println("Account: ", config.Accounts[*account].Name)
		}

		appID := config.Accounts[*account].GetAppID(config.Roles[*role])

		//Get OneLogin access Token
		token, err := onelogin.GetAccessToken(config.Onelogin.ClientID, config.Onelogin.ClientSecret)
		if err != nil {
			log.Fatalln(err)
		}

		//Get email and password from keyboard input
		var email string
		if config.Credentials.Email == "" {
			email, err = utils.PromptForString("Email")
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			email = config.Credentials.Email
		}

		var password string
		if config.Credentials.Password == "" {
			password, err = utils.PromptForSecretString("Password")
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			password = config.Credentials.Password
		}

		//SAML Assertion and MFA Devices retrieval
		assertionResponse, err := onelogin.SAMLAssertion(token, email, password, appID, config.Onelogin.AccountName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		//MFA Device verification
		var deviceID *int
		if assertionResponse.Message == onelogin.MFA_REQUIRED_STRING {
			fmt.Println("MFA Required, select a device:")
			deviceID, err = getDeviceID(assertionResponse.Devices)
			if err != nil {
				log.Fatalln(err)
			}

			var mfaCode string
			if config.Credentials.OTP == "" {
				mfaCode, err = utils.PromptForSecretString("MFA Code")
				if err != nil {
					log.Fatalln(err)
				}
			} else {
				mfaCode = config.Credentials.OTP
			}

			verificationResponse, err := onelogin.VerifyFactor(token, *deviceID, appID, assertionResponse.StateToken, mfaCode)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			assertionPayload = verificationResponse.Data
		}

		//AssumeRole With SAML on AWS
		accountID := config.Accounts[*account].AccountID
		profileName := config.Accounts[*account].ProfileName
		durationSeconds := config.Accounts[*account].DurationSeconds
		if durationSeconds == 0 {
			durationSeconds = 3600
		}
		result, err := intAWS.AssumeRoleWithSAML(accountID, config.Roles[*role], assertionPayload, durationSeconds)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = intAWS.SetCredentials(
			&intAWS.CredentialFileGetter{},
			&intAWS.CredentialFileWriter{},
			*result.Credentials.AccessKeyId,
			*result.Credentials.SecretAccessKey,
			*result.Credentials.SessionToken,
			config.DefaultRegion,
			profileName,
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Successfully set credentials for: %s\n", profileName)
	},
}

func getRole(roles []string) (*int, error) {

	roleName, err := utils.PromptSelect("Role", config.Roles, false)
	if err != nil {
		return nil, err
	}
	for k, v := range roles {
		if v == roleName {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("role not found")
}

func getAccount(accounts []Account) (*int, error) {
	var accountsName []string
	for _, v := range accounts {
		accountsName = append(accountsName, v.Name)
	}
	accountName, err := utils.PromptSelect("Account", accountsName, false)
	if err != nil {
		return nil, err
	}
	for k, v := range accounts {
		if v.Name == accountName {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("Account not found")
}

func getDeviceID(devices []onelogin.Device) (*int, error) {
	var deviceTypes []string
	for _, v := range devices {
		deviceTypes = append(deviceTypes, v.DeviceType)
	}
	selectedDeviceType, err := utils.PromptSelect("MFA Device", deviceTypes, true)
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range devices {
		if v.DeviceType == selectedDeviceType {
			return &v.DeviceID, nil
		}
	}
	return nil, fmt.Errorf("no device found")
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
