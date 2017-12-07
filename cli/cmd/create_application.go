package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/HariniGB/openstack-api/controllers"
	"fmt"
	"github.com/HariniGB/openstack-api/heat"
	"os"
)

func NewCreateApplicationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "application name [--type=...] [--flavor=...]",
		Short: "create an application",
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunCreateApplication(cmd, args); err != nil {
				log.Println(err)
			}
		},
	}
	cmd.Flags().String("type", "", "Type of application. Ex: tomcat, cirros")
	cmd.Flags().String("flavor", "", "Flavor of VM. Ex: m1.tiny, m1.small")

	return cmd
}

func RunCreateApplication(cmd *cobra.Command, args []string) error{
	name, err := NameFromCommandArgs(cmd, args)
	if err != nil {
		return err
	}

	appType := GetFlagString(cmd, "type")
	flavor := GetFlagString(cmd, "flavor")

	if appType == "" {
		return fmt.Errorf("type is a mandatory paramter")
	}

	if flavor == "" {
		return fmt.Errorf("flavor is a mandatory paramter")
	}

	user :=  os.Getenv("OS_USERNAME")
	if user == "" {
		return fmt.Errorf("source your credentials using `source ~/.virt/openrc`")
	}

	params := map[string]string{
		"name": name,
		"net_name": user,
		"subnet_name": fmt.Sprintf("%s-subnet", user),
		"flavor_name": flavor,
	}

	template := ""
	switch appType {
	case "cirros":
		template = heat.Cirros

	case "tomcat":
		template = heat.Tomcat

	default:
		return fmt.Errorf("unknow application type")
	}
	err = controllers.DeployTemplate(name, template, params)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("unable to create network")
	}
	return nil
}
