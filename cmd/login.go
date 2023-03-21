package cmd

import (
	"fmt"
	"log"
	"net/http"
	intAWS "onelogin-auth-cli/internal/aws"
	"onelogin-auth-cli/internal/onelogin"
	"onelogin-auth-cli/utils"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

type prompts struct {
}

type PromptsIface interface {
	getRole(promptRunner utils.PromptRunner, roles []string) (*int, error)
	getAccount(promptRunner utils.PromptRunner, accounts []Account) (*int, error)
	getDeviceID(promptRunner utils.PromptRunner, devices []onelogin.Device) (*int, error)
	getEmail() (string, error)
	getPassword() (string, error)
}

var httpClient = &http.Client{}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login command",
	Run: func(cmd *cobra.Command, args []string) {
		err := doLogin(onelogin.Client{}, prompts{}, intAWS.Client{}, args)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func doLogin(oneloginClient onelogin.Iface, prompts PromptsIface, awsClient intAWS.ClientIface, args []string) error {
	var role, account *int
	var err error
	var assertionPayload string

	//Get Role and Accounts from parameters or from keyboard input
	if len(args) != 2 {
		role, err = prompts.getRole(utils.GetSelectPrompt("Role", config.Roles), config.Roles)
		if err != nil {
			return err
		}
		account, err = prompts.getAccount(utils.GetSelectPrompt("Account", getAccountNames(config.Accounts)), config.Accounts)
		if err != nil {
			return err
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
	token, err := oneloginClient.GetAccessToken(httpClient, config.Onelogin.ClientID, config.Onelogin.ClientSecret)
	if err != nil {
		return err
	}

	//Get email and password from keyboard input
	var email string
	if config.Credentials.Email == "" {
		email, err = prompts.getEmail()
		if err != nil {
			return err
		}
	} else {
		email = config.Credentials.Email
	}

	var password string
	if config.Credentials.Password == "" {
		password, err = prompts.getPassword()
		if err != nil {
			return err
		}
	} else {
		password = config.Credentials.Password
	}

	//SAML Assertion and MFA Devices retrieval
	assertionResponse, err := oneloginClient.SAMLAssertion(httpClient, token, email, password, appID, config.Onelogin.AccountName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//MFA Device verification
	var deviceID *int
	if assertionResponse.Message == onelogin.MFA_REQUIRED_STRING {
		fmt.Println("MFA Required, select a device:")
		deviceID, err = prompts.getDeviceID(utils.GetSelectPrompt("MFA Device", onelogin.GetDeviceTypes(assertionResponse.Devices)), assertionResponse.Devices)
		if err != nil {
			return err
		}

		var mfaCode string
		if config.Credentials.OTP == "" {
			mfaCode, err = utils.PromptForSecretString("MFA Code")
			if err != nil {
				return err
			}
		} else {
			mfaCode = config.Credentials.OTP
		}

		verificationResponse, err := oneloginClient.VerifyFactor(httpClient, token, *deviceID, appID, assertionResponse.StateToken, mfaCode)
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
	result, err := awsClient.AssumeRoleWithSAML(accountID, config.Roles[*role], assertionPayload, durationSeconds)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = awsClient.SetCredentials(
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

	return nil
}

func (p prompts) getRole(promptRunner utils.PromptRunner, roles []string) (*int, error) {
	roleName, err := utils.PromptSelect(promptRunner, roles, false)
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

func (p prompts) getAccount(promptRunner utils.PromptRunner, accounts []Account) (*int, error) {
	accountName, err := utils.PromptSelect(promptRunner, getAccountNames(accounts), false)
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

func (p prompts) getDeviceID(promptRunner utils.PromptRunner, devices []onelogin.Device) (*int, error) {
	selectedDeviceType, err := utils.PromptSelect(promptRunner, onelogin.GetDeviceTypes(devices), true)
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

func (p prompts) getEmail() (string, error) {
	return utils.PromptForString("Email")
}
func (p prompts) getPassword() (string, error) {
	return utils.PromptForSecretString("Password")
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
