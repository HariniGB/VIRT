package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"log"
	"os"
	"io/ioutil"
)

func NewCommand() *cobra.Command {
	ensureConfigDirCreated()

	var rootCmd = &cobra.Command{
		Use:   "virt",
		Short: "Virtualization Project",
		Run: runHelp,
	}

	rootCmd.AddCommand(NewCreateCommand())
	rootCmd.AddCommand(NewLoginCommand())
	rootCmd.AddCommand(NewDeleteCommand())
	rootCmd.AddCommand(NewGetCommand())

	return rootCmd
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func NameFromCommandArgs(cmd *cobra.Command, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("NAME is required")
	}
	return args[0], nil
}

func GetFlagString(cmd *cobra.Command, flag string) string {
	s, err := cmd.Flags().GetString(flag)
	if err != nil {
		log.Fatal(fmt.Sprintf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err))
	}
	return s
}

func ensureConfigDirCreated() {
	home := os.Getenv("HOME")
	err :=os.Mkdir(fmt.Sprintf("%s/.virt", home), os.ModePerm);

	if os.IsExist(err) == false {
		log.Panic("Unable to create .virt directory due to error: ", err)
	}
}

func storeToken(content string) error {
	home := os.Getenv("HOME")
	err := ioutil.WriteFile(fmt.Sprintf("%s/.virt/openrc", home,), []byte(content), os.ModeAppend)
	if err != nil {
		return fmt.Errorf("unable to write token to file")
	}
	return nil
}

func getRcFileStr(name string, password string) string {
	return fmt.Sprintf(`
	export OS_USERNAME=%s
	export OS_PASSWORD=%s
	export OS_PROJECT_NAME=%s
	export OS_USER_DOMAIN_NAME=default
	export OS_PROJECT_DOMAIN_NAME=default
	export OS_AUTH_URL=http://10.0.0.11:5000/v3
	export OS_IDENTITY_API_VERSION=3
	export OS_IMAGE_API_VERSION=2
	`, name, password, name)
}

