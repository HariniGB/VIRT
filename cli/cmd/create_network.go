package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/HariniGB/openstack-api/controllers"
	"fmt"
	"github.com/HariniGB/openstack-api/heat"
	"os"
)

func NewCreateNetworkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "network [--cidr=...]",
		Short: "create a network topology",
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunCreateNetwork(cmd, args); err != nil {
				log.Println(err)
			}
		},
	}
	cmd.Flags().String("cidr", "", "CIDR for the subnet.")

	return cmd
}

func RunCreateNetwork(cmd *cobra.Command, args []string) error{
	cidr := GetFlagString(cmd, "cidr")

	if cidr == "" {
		return fmt.Errorf("cidr is a mandatory paramter")
	}

	user :=  os.Getenv("OS_USERNAME")
	if user == "" {
		return fmt.Errorf("source your credentials using `source ~/.virt/openrc`")
	}

	params := map[string]string{
		"net_name": user,
		"subnet_name": fmt.Sprintf("%s-subnet", user),
		"router_name": fmt.Sprintf("%s-router", user),
		"cidr": cidr,
	}
	err := controllers.DeployTemplate(fmt.Sprintf("%s-network", user), heat.Network, params)
	if err != nil {
		return fmt.Errorf("unable to create network")
	}
	return nil
}
