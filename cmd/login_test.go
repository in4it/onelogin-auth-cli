package cmd

import (
	"onelogin-auth-cli/internal/onelogin"
	"testing"
)

type promptMock struct {
	RunResultInt    int
	RunResultString string
}

func (p promptMock) Run() (int, string, error) {
	return p.RunResultInt, p.RunResultString, nil
}

func TestGetRole(t *testing.T) {
	options := []string{"a", "b", "c"}
	answer := 2
	res, err := getRole(promptMock{RunResultString: options[answer]}, options)
	if err != nil {
		t.Fatalf("getRole error: %s", err)
	}
	if *res != answer {
		t.Fatalf("Got wrong answer: %d", *res)
	}
}

func TestGetRoleEmpty(t *testing.T) {
	options := []string{}
	_, err := getRole(promptMock{RunResultString: ""}, options)
	if err.Error() != "role not found" {
		t.Fatalf("Expected role not found error, got nil")
	}
}

func TestGetRoleOneItem(t *testing.T) {
	options := []string{"a"}
	answer := 0
	res, err := getRole(promptMock{RunResultString: options[answer]}, options)
	if err != nil {
		t.Fatalf("getRole error: %s", err)
	}
	if *res != answer {
		t.Fatalf("Got wrong answer: %d", *res)
	}
}
func TestGetRoleWrongItem(t *testing.T) {
	options := []string{"a", "b", "c"}
	_, err := getRole(promptMock{RunResultString: "d"}, options)
	if err.Error() != "role not found" {
		t.Fatalf("Expected role not found error, got nil")
	}
}

func TestGetAccount(t *testing.T) {
	accounts := []Account{
		{
			Name: "a",
		},
		{
			Name: "b",
		},
		{
			Name: "c",
		},
	}
	answer := 1
	res, err := getAccount(promptMock{RunResultString: accounts[answer].Name}, accounts)
	if err != nil {
		t.Fatalf("getAccount error: %s", err)
	}
	if *res != answer {
		t.Fatalf("Got wrong answer: %d", *res)
	}
}

func TestGetDeviceID(t *testing.T) {
	devices := []onelogin.Device{
		{
			DeviceID:   1,
			DeviceType: "type1",
		},
		{
			DeviceID:   2,
			DeviceType: "type2",
		},
		{
			DeviceID:   3,
			DeviceType: "type3",
		},
	}
	answer := 2
	res, err := getDeviceID(promptMock{RunResultString: devices[answer].DeviceType}, devices)
	if err != nil {
		t.Fatalf("getDeviceID error: %s", err)
	}
	if *res != devices[answer].DeviceID {
		t.Fatalf("Got wrong answer: %d", *res)
	}
}
