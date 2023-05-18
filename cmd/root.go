package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Onelogin      OneLoginConf
	Credentials   Credentials `yaml:"credentials"`
	Accounts      []Account   `yaml:"accounts"`
	Roles         []string    `yaml:"roles"`
	DefaultRegion string      `yaml:"defaultRegion"`
}

type Credentials struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	OTP      string `yaml:"otp"`
}

type OneLoginConf struct {
	ClientID     string `yaml:"onelogin-client-id"`
	ClientSecret string `yaml:"onelogin-client-secret"`
	AccountName  string `yaml:"onelogin-account"`
}
type Account struct {
	Name            string            `yaml:"name"`
	AppID           string            `yaml:"appID"`
	AppIDsByRole    map[string]string `yaml:"appIDsByRole"`
	AccountID       string            `yaml:"accountID"`
	ProfileName     string            `yaml:"profileName"`
	DurationSeconds int64             `yaml:"durationSeconds"`
}

func getAccountNames(accounts []Account) []string {
	var accountsName []string
	for _, v := range accounts {
		accountsName = append(accountsName, v.Name)
	}
	return accountsName
}

func (a *Account) GetAppID(role string) string {
	if a.AppIDsByRole != nil {
		if appID, ok := a.AppIDsByRole[role]; ok {
			return appID
		}
	}
	return a.AppID
}

var version string
var config Config
var configFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "onelogin-auth",
	Short: "OneLogin authentication CLI",
}

func SetVersion(v string) {
	version = v
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(LoadConfig)

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is config.yaml)")
	viper.BindPFlag("config", rootCmd.Flags().Lookup("config"))
}

func LoadConfig() {
	userDefinedConfigFile := os.Getenv("ONELOGIN_AUTH_CLI_CONFIG_FILE")
	if userDefinedConfigFile != "" {
		viper.SetConfigFile(userDefinedConfigFile)
	} else {
		if configFile == "" {
			viper.AddConfigPath("./")
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
		} else {
			viper.SetConfigFile(configFile)
		}
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	viper.BindEnv("credentials.email", "EMAIL")
	viper.BindEnv("credentials.password", "PASSWORD")
	viper.BindEnv("credentials.otp", "OTP")

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}
}
