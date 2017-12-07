package cmd

import (
	"github.com/spf13/cobra"
	"github.com/prometheus/log"
	"github.com/HariniGB/openstack-api/controllers"
	"fmt"
)

func NewCreateAccountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account username [--password=...]",
		Short: "create an account",
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunCreateAccount(cmd, args); err != nil {
				log.Println(err)
			}
		},
	}
	cmd.Flags().String("password", "", "Passowrd to authenticate with.")

	return cmd
}

func RunCreateAccount(cmd *cobra.Command, args []string) error{
	name, err := NameFromCommandArgs(cmd, args)
	if err != nil {
		return err
	}

	password := GetFlagString(cmd, "password")

	err = controllers.CreateAccount(name, password)
	if err != nil {
		return fmt.Errorf("unable to create user")
	}

	rc := getRcFileStr(name, password)
	err = storeToken(rc)
	if err != nil {
		return err
	}
	return nil
}
