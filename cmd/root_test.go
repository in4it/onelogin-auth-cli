package cmd

import (
	"os"
	"testing"
)

func TestGetAppIDNonExistent(t *testing.T) {
	a := Account{}
	appID := a.GetAppID("role-name")
	if appID != "" {
		t.Fatalf("Unexpected app ID: %s", appID)
	}
}

func TestGetAppID(t *testing.T) {
	a := Account{
		AppIDsByRole: map[string]string{
			"test-role": "123456",
		},
	}
	appID := a.GetAppID("test-role")
	if appID != "123456" {
		t.Fatalf("Unexpected app ID: %s", appID)
	}
}

func TestLoadConfigEnvVar(t *testing.T) {
	os.Setenv("ONELOGIN_AUTH_CLI_CONFIG_FILE", "../internal/testdata/config.yaml")
	LoadConfig()

	if config.Onelogin.AccountName != "testdata" {
		t.Fatalf("config variable doesn't contain testdata")
	}

}
func TestLoadConfig(t *testing.T) {
	os.Setenv("ONELOGIN_AUTH_CLI_CONFIG_FILE", "")
	configFile = "../internal/testdata/config.yaml"
	LoadConfig()

	if config.Onelogin.AccountName != "testdata" {
		t.Fatalf("config variable doesn't contain testdata")
	}

}
