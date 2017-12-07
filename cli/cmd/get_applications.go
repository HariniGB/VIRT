package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/HariniGB/openstack-api/controllers"
	"fmt"
)

func NewGetApplicationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "applications",
		Short: "get all applications",
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunGetApplications(cmd, args); err != nil {
				log.Println(err)
			}
		},
	}
	return cmd
}

func RunGetApplications(cmd *cobra.Command, args []string) error{

	stacks := controllers.GetStacks()
	if stacks == nil {
		return fmt.Errorf("unable to get applications")
	}

	rows := [][]string{}
	for _, stack := range stacks {
		row := []string{stack.Name, stack.Status, stack.PrivateIp, stack.PublicIp, stack.CreationTime}
		rows = append(rows, row)
	}

	newTable(rows)
	return nil
}
