package cmd

import (
	"github.com/spf13/cobra"
)

func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get an object",
		Run: runHelp,
	}

	cmd.AddCommand(NewGetApplicationCommand())
	cmd.AddCommand(NewGetApplicationsCommand())
	return cmd
}
