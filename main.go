package main

import (
	"onelogin-auth-cli/cmd"
)

// GoReleaser will set this value at build time
var version = "development"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
