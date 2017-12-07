package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/HariniGB/openstack-api/controllers"
	"fmt"
)

func NewDeleteApplicationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete name",
		Short: "delete an application",
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunDeleteApplication(cmd, args); err != nil {
				log.Println(err)
			}
		},
	}
	return cmd
}

func RunDeleteApplication(cmd *cobra.Command, args []string) error{
	name, err := NameFromCommandArgs(cmd, args)
	if err != nil {
		return err
	}

	err = controllers.UndeployTemplate(name)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("unable to delete application")
	}

	fmt.Printf("%s is successfully deleted", name)
	return nil
}
