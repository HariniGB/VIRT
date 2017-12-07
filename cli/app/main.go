package app

import "github.com/HariniGB/openstack-api/cli/cmd"

func Run() error {
	c := cmd.NewCommand()
	return c.Execute()
}

