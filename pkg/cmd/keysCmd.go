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
	Short: "Manage SSH keys for secure server access",
	Long: `The "keys" command group provides options for managing SSH keys, 
including deploying and authorizing keys for secure server access.

To remove an authorized key, navigate to /home/mole/.ssh and delete the 
relevant entry in the authorized_keys file.

To regenerate the deploy key, navigate to /home/mole/.ssh and delete 
the id_rsa and id_rsa.pub files. This will allow the system to create 
a new deploy key as needed.`,
}

var readDeployKeyCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Retrieve or create the deploy key for SSH access",
	Long: `The "deploy" command displays the current deploy key (id_rsa.pub),
used for accessing private repositories securely. 

If no deploy key is found, a new one will be generated automatically 
and saved to the standard SSH key path. This deploy key enables 
secure, automated interactions with external repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
		key, err := actions.FindOrCreateDeployKey()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(key)
	},
}

var addAuthorizedKeyCmd = &cobra.Command{
	Use:   "authorize [public RSA key]",
	Short: "Add a new public key to the authorized_keys file",
	Long: `The "authorize" command validates and appends a given public RSA key 
to the authorized_keys file. This allows the specified key to be used 
for SSH access to the server.

Ensure the key provided is correctly formatted, as it will be validated 
before being added to prevent errors. Only unique keys will be appended 
to avoid duplicates in the authorized_keys file.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := actions.AddAuthorizedKeys(strings.Join(args, " "))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Key added")
	},
}
