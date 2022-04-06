package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Onelogin      OneLoginConf
	Accounts      []Account `yaml:"accounts"`
	Roles         []string  `yaml:"roles"`
	DefaultRegion string    `yaml:"defaultRegion"`
}

type OneLoginConf struct {
	ClientID     string `yaml:"onelogin-client-id"`
	ClientSecret string `yaml:"onelogin-client-secret"`
	AccountName  string `yaml:"onelogin-account"`
}
type Account struct {
	Name        string `yaml:"name"`
	AppID       string `yaml:"appID"`
	AccountID   string `yaml:"accountID"`
	ProfileName string `yaml:"profileName"`
}

var config Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "onelogin-auth",
	Short: "OneLogin authenticatio CLI Tool",
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

	err = viper.Unmarshal(&config)
	return
}
