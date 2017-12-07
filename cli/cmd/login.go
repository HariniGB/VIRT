package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/HariniGB/openstack-api/controllers"
	"fmt"
	"gopkg.in/bufio.v1"
	"os"
	"strings"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

func NewLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "login [--username=...] [--password=...]",
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunLogin(cmd, args); err != nil {
				log.Println(err)
			}
		},
	}
	cmd.Flags().String("username", "", "Username to authenticate with.")
	cmd.Flags().String("password", "", "Passowrd to authenticate with.")

	return cmd
}

func RunLogin(cmd *cobra.Command, args []string) error{
	username := GetFlagString(cmd, "username")
	password := GetFlagString(cmd, "password")

	if username == "" {
		fmt.Print("Username: ")
		reader := bufio.NewReader(os.Stdin)
		var err error
		username, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("unable to get username: %s", err)
		}
		username = strings.TrimSpace(username)
	}

	if password == "" {
		fmt.Print("Username: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("unable to get user password: %s", err)
		}
		password = string(bytePassword)
		password = strings.TrimSpace(password)
	}

	_, err := controllers.KeystoneLogin(username, password, "")
	if err != nil {
		return fmt.Errorf("unable to create user")
	}

	rc := getRcFileStr(username, password)
	err = storeToken(rc)
	if err != nil {
		return err
	}

	fmt.Println("Successfully authenticated. Please run `source ~/.virt/openrc` before running other commands")
	return nil
}