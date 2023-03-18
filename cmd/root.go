package cmd

import (
	"log"

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
	Name            string `yaml:"name"`
	AppID           string `yaml:"appID"`
	AccountID       string `yaml:"accountID"`
	ProfileName     string `yaml:"profileName"`
	DurationSeconds int64  `yaml:"durationSeconds"`
}

var version string
var config Config

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
	var err error
	config, err = LoadConfig("./")
	if err != nil {
		log.Fatalln(err)
	}
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	viper.BindEnv("credentials.email", "EMAIL")
	viper.BindEnv("credentials.password", "PASSWORD")
	viper.BindEnv("credentials.otp", "OTP")

	err = viper.Unmarshal(&config)
	return
}
