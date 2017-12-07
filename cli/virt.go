package main

import (
	"os"

	"github.com/HariniGB/openstack-api/cli/app"
)

func main() {
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
