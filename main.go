package main

import (
	"fmt"
	"os"

	"github.com/tzapu/disco-bit/cmd"
)

var (
	buildVersion string
)

func main() {
	// proxy for version and date
	cmd.BuildVersion = buildVersion
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
