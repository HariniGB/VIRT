package cmd

import (
	"github.com/spf13/cobra"
)

func NewDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete an object",
		Run: runHelp,
	}

	cmd.AddCommand(NewDeleteApplicationCommand())
	return cmd
}
