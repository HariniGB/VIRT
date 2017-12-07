package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/HariniGB/openstack-api/controllers"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
)

var (
	table_headers = []string{"Name", "Status", "Private Ip", "Public Ip", "Creation Time"}
	newTable = func(data [][]string) *tablewriter.Table {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(table_headers)
		table.AppendBulk(data)
		table.Render()
		return table
	}
)

func NewGetApplicationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "application name",
		Short: "get an application",
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunGetApplication(cmd, args); err != nil {
				log.Println(err)
			}
		},
	}
	return cmd
}

func RunGetApplication(cmd *cobra.Command, args []string) error{
	name, err := NameFromCommandArgs(cmd, args)
	if err != nil {
		return err
	}

	stack := controllers.GetStack(name)
	if stack == nil {
		return fmt.Errorf("unable to get application")
	}

	row := []string{stack.Name, stack.Status, stack.PrivateIp, stack.PublicIp, stack.CreationTime}
	newTable([][]string{row})
	return nil
}
