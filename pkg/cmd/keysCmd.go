package cmd

import (
	"fmt"
	"strings"

	"github.com/6oof/mole/pkg/actions"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(keysRootCmd)

	keysRootCmd.AddCommand(readDeployKeyCmd)
	keysRootCmd.AddCommand(addAuthorizedKeyCmd)
}

var keysRootCmd = &cobra.Command{
	Use:   "keys",
	Short: "Interact with ssh keys",
	Long: `Keys is a group of commands for adding ssh keys.

to remove an authorized key you should cd into /home/mole/.ssh
and deleta an entry from authorized_keys.

To regenerate the deploy key, cd into /home/home/mole/.ssh
and delete (rm) id_rsa and id_rsa.pub`,
}

var readDeployKeyCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy returns the deploy key",
	Long: `Deploy returns the deploy key (id_rsa.pub).

This key can be used to interract with private repositories.

If the deploy key is not found, it will be created.`,
	Run: func(cmd *cobra.Command, args []string) {
		key, err := actions.FindOrCreateDeployKey()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(key)
	},
}

var addAuthorizedKeyCmd = &cobra.Command{
	Use:   "authorize [public RSA key]",
	Short: "Authorize adds an authorized key",
	Long: `Authorize adds an authorized key to the authorized_keys file.

The key is validated before it is added.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := actions.AddAuthorizedKeys(strings.Join(args, " "))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Key added")
		}
	},
}
