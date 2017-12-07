package cmd

import (
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create an object",
		Run: runHelp,
	}

	cmd.AddCommand(NewCreateAccountCommand())
	cmd.AddCommand(NewCreateNetworkCommand())
	cmd.AddCommand(NewCreateApplicationCommand())

	return cmd
}
