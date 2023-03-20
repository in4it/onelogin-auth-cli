package main

import (
	"log"
	"onelogin-auth-cli/cmd"
)

// GoReleaser will set this value at build time
var version = "development"

func main() {
	cmd.SetVersion(version)
	var err error
	err = cmd.LoadConfig("./")
	if err != nil {
		log.Fatalln(err)
	}
	cmd.Execute()
}
